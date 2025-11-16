package cmd

import (
	"fmt"
	"os"
	"path"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func New() error {
	cmd := &cobra.Command{
		Use:   "neomux",
		Short: "neovim multiplexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	{
		sc, err := newNewCmd()
		if err != nil {
			return fmt.Errorf("new \"new\" cmd: %w", err)
		}
		cmd.AddCommand(sc)
	}
	cmd.AddCommand(newNvCmd())
	cmd.AddCommand(newKillCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newStateCmd())
	cmd.AddCommand(newDuplicateCmd())
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

func newState() (adapter.Adapter, error) {
	path, err := statePath()
	if err != nil {
		return nil, fmt.Errorf("state path: %w", err)
	}
	state, err := adapter.New(path)
	if err != nil {
		return nil, fmt.Errorf("sqlite state adapter: %w", err)
	}
	return state, nil
}
