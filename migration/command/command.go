package migration

import (
	"fmt"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/database/database_client"
	"frisboo-bank/pkg/database/migration"
	"frisboo-bank/pkg/environment"
	"os"

	migrationcommandtype "frisboo-bank/pkg/database/migration/contracts/enums/migration_command_type"

	"github.com/spf13/cobra"
)

type MigrationCommand struct{}

var rootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run the db migrations",
}

var cmdUp = &cobra.Command{
	Use:   migrationcommandtype.MigrationCommandTypes.UP.String(),
	Short: "Run an up migration",
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(cmd, migrationcommandtype.MigrationCommandTypes.UP)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var cmdDown = &cobra.Command{
	Use:   migrationcommandtype.MigrationCommandTypes.DOWN.String(),
	Short: "Run a down migration",
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(cmd, migrationcommandtype.MigrationCommandTypes.DOWN)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func NewMigrationCommand() *MigrationCommand {
	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdDown)

	return &MigrationCommand{}
}

func (c *MigrationCommand) Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}

	fmt.Println("Migration completed...")

	return nil
}

func executeMigration(cmd *cobra.Command, commandType migrationcommandtype.MigrationCommandType) error {
	fmt.Println("Migration process started...")

	app := module.ModuleFunc(
		"migration-app",
		environment.ModuleFunc(environment.Development),
		database_client.Module,
		migration.Module,
	)

	// config := &config.MigrationOptions{
	// 	Host:         "localhost",
	// 	Port:         "5432",
	// 	User:         "postgres",
	// 	DBName:       "customers-service",
	// 	SSLMode:      false,
	// 	Password:     "postgres",
	// 	MigrationDir: "db/migrations",
	// }
	//
	// migrationRunner, err := gomigrate.NewPostgresqlMigrator(config)
	// if err != nil {
	// 	return err
	// }
	//
	// switch commandType {
	// case config.CommandTypeUp:
	// 	return migrationRunner.Up(context.Background(), 0)
	// case config.CommandTypeDown:
	// 	return migrationRunner.Down(context.Background(), 0)
	// }
	//
	// panic(syserrors.Newf("migration: unsupported command type: `%s`", commandType))
}
