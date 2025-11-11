package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
)

func New() error {
	cmd := &cobra.Command{
		Use:   "neomux",
		Short: "neovim multiplexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newNewCmd())
	cmd.AddCommand(newNvCmd())
	cmd.AddCommand(newKillCmd())
	cmd.AddCommand(newListCmd())
	return cmd.Execute()
}

func statePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home dir: %w", err)
	}
	path := path.Join(home, ".cache", "nvim", "servers.db")
	return path, nil
}
