package command

import (
	"fmt"

	"frisboo-bank/pkg/application/builder"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/migration"

	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"

	"github.com/spf13/cobra"
)

type MigrationCommand struct{}

var rootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run the db migrations",
}

var cmdUp = &cobra.Command{
	Use:   fmt.Sprintf("%s [database-name]", migrationcommandtype.MigrationCommandTypes.UP.String()),
	Short: "Run an up migration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(migrationcommandtype.MigrationCommandTypes.UP, cmd, args)
		if err != nil {
			panic(err)
		}
	},
}

var cmdDown = &cobra.Command{
	Use:   fmt.Sprintf("%s [database-name]", migrationcommandtype.MigrationCommandTypes.DOWN.String()),
	Short: "Run a down migration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(migrationcommandtype.MigrationCommandTypes.DOWN, cmd, args)
		if err != nil {
			panic(err)
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
	return nil
}

func executeMigration(commandType migrationcommandtype.MigrationCommandType, cmd *cobra.Command, args []string) error {
	databaseName := args[0]

	fmt.Printf("Migration started for database %s...\n", databaseName)

	appBuilder, err := builder.NewApplicationBuilder(environment.Environments.DEVELOPMENT)
	if err != nil {
		return err
	}

	m := module.ModuleFunc("migration",
		migration.ModuleFunc(appBuilder),
	)
	appBuilder.ProvideModule(m)

	app := appBuilder.Build()

	if err := app.Run(); err != nil {
		return err
	}

	fmt.Printf("Migration of the database %s done successfully\n", databaseName)

	return nil
}
