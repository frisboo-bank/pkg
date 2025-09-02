package application

import (
	"fmt"

	"github.com/spf13/cobra"
)

type ApplicationCliConfig struct {
	Name        string
	Header      func()
	Description string
}

type ApplicationCli struct {
	cfg *ApplicationCliConfig
}

func NewApplicationCli(cfg *ApplicationCliConfig) *ApplicationCli {
	return &ApplicationCli{
		cfg: cfg,
	}
}

func (c *ApplicationCli) Execute(bootstrapFn func(cmd *cobra.Command, args []string)) error {
	fmt.Println("")
	if c.cfg.Header != nil {
		c.cfg.Header()
	}

	rootCmd := cobra.Command{
		Use:   c.cfg.Name,
		Short: c.cfg.Description,
		Run:   bootstrapFn,
	}

	return rootCmd.Execute()
}
