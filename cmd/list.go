package cmd

import (
	"context"
	"fmt"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all nvim servers' labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := statePath()
			if err != nil {
				return fmt.Errorf("state path: %w", err)
			}
			state, err := adapter.New(path)
			if err != nil {
				return fmt.Errorf("new sqlite adapter: %w (path=%q)", err, path)
			}
			servers, err := state.ListServers(context.Background())
			if err != nil {
				return fmt.Errorf("list servers: %w", err)
			}
			for _, s := range servers {
				fmt.Println(s.Label)
			}
			return nil
		},
	}
}
