package migration

import (
	"context"
	"errors"
	"fmt"
	"os"

	"frisboo-bank/pkg/db/migration/gomigrate"

	migrationConfig "frisboo-bank/pkg/db/migration/config"

	"github.com/spf13/cobra"
)

func init() {
}

type MigrationCommand struct{}

var rootCmd = &cobra.Command{
	Use:   "migration",
	Short: "Run the db migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return nil
		}

		cmd.SetArgs([]string{"up"})
		return cmd.Execute()
	},
}

var cmdUp = &cobra.Command{
	Use:   "up",
	Short: "Run an up migration",
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(cmd, migrationConfig.CommandTypeUp)
		if err != nil {
			fmt.Println(fmt.Errorf(migrationConfig.ErrMigrationFailed.Error(), err))
			os.Exit(1)
		}
	},
}

var cmdDown = &cobra.Command{
	Use:   "down",
	Short: "Run a down migration",
	Run: func(cmd *cobra.Command, args []string) {
		err := executeMigration(cmd, migrationConfig.CommandTypeDown)
		if err != nil {
			fmt.Println(errors.Join(migrationConfig.ErrMigrationFailed, err))
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

func executeMigration(cmd *cobra.Command, commandType migrationConfig.CommandType) error {
	fmt.Println("Migration process started...")

	config := &migrationConfig.MigrationOptions{
		Host:         "localhost",
		Port:         "5432",
		User:         "postgres",
		DBName:       "customers-service",
		SSLMode:      false,
		Password:     "postgres",
		MigrationDir: "db/migrations",
	}

	migrationRunner, err := gomigrate.NewPostgresqlMigrator(config)
	if err != nil {
		return err
	}

	switch commandType {
	case migrationConfig.CommandTypeUp:
		return migrationRunner.Up(context.Background(), 0)
	case migrationConfig.CommandTypeDown:
		return migrationRunner.Down(context.Background(), 0)
	}

	panic(fmt.Errorf("migration: unsupported command type: `%s`", commandType))
}
