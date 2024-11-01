package mariadb_dump

import (
	"bytes"
	"fmt"
	"github.com/denisakp/sentinel/internal/backup/sql"
	"github.com/denisakp/sentinel/internal/storage"
	"github.com/denisakp/sentinel/internal/utils"
	"os/exec"
)

func Backup(mda *MariaDBDumpArgs) error {
	// Validate the required arguments
	args, err := ArgsBuilder(mda)
	if err != nil {
		return fmt.Errorf("failed to build arguments: %w", err)
	}

	// check connectivity
	if ok, err := sql.CheckConnectivity("mysql", mda.Host, mda.Port, mda.Username, mda.Password, mda.Database); !ok {
		return err
	}

	// execute mariadb-dump command
	cmd := exec.Command("mariadb-dump", args...)

	// capture command error
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// capture command output
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute maridb-dump command - %w, %s", err, stdErr.String())
	}

	// get the storage handler
	storageHandler, err := storage.NewStorage(mda.StorageType)
	if err != nil {
		return err
	}

	// get the backup path
	backupPath, err := storageHandler.GetBackupPath(mda.StoragePath)

	// set output name with customizable extension (default is .sql)
	mda.OutName = utils.FinalOutName(mda.OutName)

	// get the full path
	fullPath := utils.FullPath(backupPath, mda.OutName)

	if err := storageHandler.WriteBackup(stdOut.Bytes(), fullPath); err != nil {
		return err
	}

	fmt.Printf("Backup file created at %s\n", fullPath)

	return nil
}