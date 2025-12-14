package cmd

import (
	"fmt"
	"os"
	"path"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func New() error {
	var path *string
	cmd := &cobra.Command{
		Use:   "neomux",
		Short: "neovim multiplexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	defaultPath, err := defaultStatePath()
	if err != nil {
		return fmt.Errorf("defaultStatePath: %w", err)
	}
	path = cmd.PersistentFlags().String("state", defaultPath, "state path")
	state, err := newState(*path)
	if err != nil {
		return fmt.Errorf("new state: %w", err)
	}
	{
		sc, err := newNewCmd(state)
		if err != nil {
			return fmt.Errorf("new \"new\" cmd: %w", err)
		}
		cmd.AddCommand(sc)
	}
	cmd.AddCommand(newNvCmd(state))
	cmd.AddCommand(newKillCmd(state))
	cmd.AddCommand(newListCmd(state))
	cmd.AddCommand(newStateCmd(state))
	cmd.AddCommand(newDuplicateCmd(state))
	return cmd.Execute()
}

func defaultStatePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("user home dir: %w", err)
	}
	path := path.Join(home, ".cache", "nvim", "servers.db")
	return path, nil
}

func newState(path string) (adapter.Adapter, error) {
	state, err := adapter.New(path)
	if err != nil {
		return nil, fmt.Errorf("sqlite state adapter: %w", err)
	}
	return state, nil
}
