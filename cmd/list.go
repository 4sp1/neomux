package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all nvim servers' labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}
			servers, err := state.ListServers(context.Background())
			if err != nil {
				return fmt.Errorf("list servers: %w", err)
			}
			for _, s := range servers {
				fmt.Printf("%s:%s\n", s.Label, s.Workdir)
			}
			return nil
		},
	}
}
