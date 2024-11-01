package pg_dump

import (
	"bytes"
	"fmt"
	"github.com/denisakp/sentinel/internal/backup/sql"
	"github.com/denisakp/sentinel/internal/storage"
	"github.com/denisakp/sentinel/internal/utils"
	"os/exec"
)

// Backup backs up a PostgresSQL database using pg_dump
func Backup(pda *PgDumpArgs) error {
	// get the storage handler
	storageHandler, err := storage.NewStorage(pda.StorageType)
	if err != nil {
		return err
	}

	backupPath, err := storageHandler.GetBackupPath(pda.StoragePath)
	if err != nil {
		return err
	}

	// build pg_dump arguments
	args, err := argsBuilder(pda)
	if err != nil {
		return fmt.Errorf("failed to build pg_dump args - %w", err)
	}

	if err := validateRequiredArgs(pda); err != nil {
		return err
	}

	// check connectivity
	if ok, err := sql.CheckConnectivity("postgres", pda.Host, pda.Port, pda.Username, pda.Password, pda.Database); !ok {
		return err
	}

	// set the output name
	pda.OutName = utils.FullPath(backupPath, pda.OutName)

	// run pg_dump command
	cmd := exec.Command("pg_dump", args...)

	// capture the command error
	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	// capture the command output
	var stdOut bytes.Buffer
	cmd.Stdout = &stdOut

	// remove the password from the environment after the command is done
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", pda.Password)) // set the password in the environment
	defer func() {
		cmd.Env = cmd.Env[:len(cmd.Env)-1]
	}()

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute pg_dump command - %w, %s", err, stdErr.String())
	}

	// write the backup to the storage
	if err := storageHandler.WriteBackup(stdOut.Bytes(), pda.OutName); err != nil {
		return err
	}

	fmt.Printf("Backup file created at %s\n", pda.OutName)

	return nil
}