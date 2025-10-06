package command

import (
	"fmt"

	"frisboo-bank/pkg/application/builder"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	databaseclient "frisboo-bank/pkg/database/database_client"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/migration"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
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
	cmdUp.Flags().Uint("version", 0, "Migration version")
	cmdDown.Flags().Uint("version", 0, "Migration version")

	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdDown)

	return &MigrationCommand{}
}

func (c *MigrationCommand) Execute() error {
	return rootCmd.Execute()
}

func executeMigration(commandType migrationcommandtype.MigrationCommandType, cmd *cobra.Command, args []string) error {
	databaseName := args[0]
	version, err := cmd.Flags().GetUint("version")
	if err != nil {
		return err
	}

	fmt.Printf("Migration started for database %s...\n", databaseName)

	appBuilder, err := builder.NewApplicationBuilder(environment.Environments.DEVELOPMENT)
	if err != nil {
		return err
	}

	cfgRegistry, err := config.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
	if err != nil {
		return fmt.Errorf("failed to load migration registry: %w", err)
	}

	cfg, err := cfgRegistry.GetByName(databaseName)
	if err != nil {
		return err
	}

	m := module.ModuleFunc(
		"migration",
		databaseclient.ModuleFunc(appBuilder),
		migration.ModuleFunc(migration.ModuleProps{
			AppBuilder:  appBuilder,
			CfgRegistry: cfgRegistry,
		}),
	)

	m.AddInvoker(invoker.InvokerFunc(func(props struct {
		Migrator contracts.Migrator `name:"migrationRef"`
	},
	) {
		migrator := props.Migrator

		migrator.Logger().Info("Migration process started...")

		var err error
		switch commandType {
		case migrationcommandtype.MigrationCommandTypes.UP:
			err = migrator.Up(version)
		case migrationcommandtype.MigrationCommandTypes.DOWN:
			err = migrator.Down(version)
		}
		if err != nil {
			migrator.Logger().Fatalf("migration failed with error: %v", err)
		}

		migrator.Logger().Info("Migration completed...")
	},
		invoker.NamedDep("migrationRef", fmt.Sprintf(migration.MigrationsProvider, cfg.DB)),
	))

	appBuilder.ProvideModule(m)

	app := appBuilder.Build()

	if err := app.Run(); err != nil {
		return err
	}

	fmt.Printf("Migration of the database %s done successfully\n", databaseName)
	return nil
}
