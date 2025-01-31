package mariadb_dump

import (
	"fmt"
	"github.com/denisakp/sentinel/internal/backup"
	"github.com/denisakp/sentinel/internal/storage"
	"github.com/denisakp/sentinel/internal/utils"
)

type MariaDBDumpArgs struct {
	Host           string          // MariaDB host
	Port           string          // MariaDB port
	Username       string          // MariaDB username
	Password       string          // MariaDB password
	Database       string          // MariaDB database name
	AdditionalArgs string          // Additional arguments for the mariadb_dump command
	Storage        *storage.Params // Storage parameters
}

// ArgsBuilder builds the arguments for the mariadb_dump command
func ArgsBuilder(mda *MariaDBDumpArgs) ([]string, error) {
	if err := validateRequiredArgs(mda); err != nil {
		return nil, err
	}

	// set the default host and port if not provided
	mda.Host = utils.DefaultValue(mda.Host, "127.0.0.1")
	mda.Port = utils.DefaultValue(mda.Port, "3306")

	// build the required arguments
	args := []string{
		fmt.Sprintf("--host=%s", mda.Host),
		fmt.Sprintf("--port=%s", mda.Port),
		fmt.Sprintf("--user=%s", mda.Username),
	}

	if mda.Password != "" {
		args = append(args, fmt.Sprintf("--password=%s", mda.Password))
	} // add the password argument if provided

	if mda.AdditionalArgs != "" {
		additionalArgs := backup.ParseAdditionalArgs(mda.AdditionalArgs)
		args = append(args, additionalArgs...)
	} // add additional arguments if provided

	args = backup.RemoveArgsDuplicate(args) // remove duplicated arguments
	args = append(args, mda.Database)       // add the database name to the arguments

	return args, nil
}
