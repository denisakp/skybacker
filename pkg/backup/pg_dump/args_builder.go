package pg_dump

import (
	"fmt"
	"github.com/denisakp/sentinel/internal/backup"
	"github.com/denisakp/sentinel/internal/storage"
	"github.com/denisakp/sentinel/internal/utils"
)

type PgDumpArgs struct {
	Host                 string          // PostgresSQL host
	Port                 string          // PostgresSQL port
	Username             string          // PostgresSQL username
	Password             string          // PostgresSQL password
	Database             string          // PostgresSQL database name
	PgOutFormat          string          // Output format for the backup file
	Compress             bool            // Enable compression
	CompressionAlgorithm string          // Compression algorithm
	CompressionLevel     int             // Compression level
	AdditionalArgs       string          // Additional arguments for the pg_dump command
	Storage              *storage.Params // Storage parameters
}

// argsBuilder builds the arguments for the pg_dump command
func argsBuilder(pda *PgDumpArgs, backupPath string) ([]string, error) {
	if err := validateRequiredArgs(pda); err != nil {
		return nil, err
	}

	// initialize default arguments
	initializeDefaultArgs(pda)

	if err := validatePgOutFormat(pda.PgOutFormat); err != nil {
		return nil, err
	}

	// validate output format
	if err := validatePgOutFormat(pda.PgOutFormat); err != nil {
		return nil, err
	}

	// handle backup outName
	if err := setOutName(pda); err != nil {
		return nil, err
	}

	args := []string{
		fmt.Sprintf("--host=%s", pda.Host),
		fmt.Sprintf("--port=%s", pda.Port),
		fmt.Sprintf("--username=%s", pda.Username),
		fmt.Sprintf("--dbname=%s", pda.Database),
		fmt.Sprintf("--format=%s", pda.PgOutFormat),
	}

	if pda.Compress {
		if err := addCompression(&args, pda); err != nil {
			return nil, err
		}
	}

	if pda.PgOutFormat == "d" {
		if pda.Storage.StorageType != "local" {
			pda.Storage.OutName = utils.FormatResourceValue(pda.Storage.OutName)
		} else {
			pda.Storage.OutName = utils.FullPath(backupPath, pda.Storage.OutName)
		}
		args = append(args, fmt.Sprintf("--file=%s", pda.Storage.OutName))
	}

	// handle additional arguments
	if pda.AdditionalArgs != "" {
		additionalArgs := backup.ParseAdditionalArgs(pda.AdditionalArgs)
		args = append(args, additionalArgs...)
	}

	// remove duplicated arguments
	args = backup.RemoveArgsDuplicate(args) // remove duplicated arguments

	return args, nil
}

func addCompression(args *[]string, pda *PgDumpArgs) error {
	// set the default compression algorithm to gzip if not provided
	pda.CompressionAlgorithm = utils.DefaultValue(pda.CompressionAlgorithm, "gzip")
	if err := validatePgCompressionAlgorithm(pda.CompressionAlgorithm); err != nil {
		return err
	}

	// validate the compression level
	if err := validatePgCompressionLevel(pda.CompressionLevel); err != nil {
		return err
	}

	// add the compression arguments
	*args = append(*args, fmt.Sprintf("--compress=%s:%d", pda.CompressionAlgorithm, pda.CompressionLevel))

	return nil
}

func initializeDefaultArgs(pda *PgDumpArgs) {
	pda.Host = utils.DefaultValue(pda.Host, "127.0.0.1")
	pda.Port = utils.DefaultValue(pda.Port, "5432")
	pda.PgOutFormat = utils.DefaultValue(pda.PgOutFormat, "p")

	if pda.CompressionAlgorithm != "" {
		pda.Compress = true
	}
}
