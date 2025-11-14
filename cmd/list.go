package cmd

import (
	"context"
	"fmt"
	"strings"

	procs "github.com/4sp1/neomux/internal/adapter/os/process"
	"github.com/4sp1/neomux/internal/app"
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

			proc, err := procs.New()
			if err != nil {
				return fmt.Errorf("procs adapter: new: %w", err)
			}

			app, err := app.New(proc, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			labels, err := app.StateClean()
			if err != nil {
				return fmt.Errorf("app: state clean: %w", err)
			}
			{
				ls := make([]string, len(labels))
				for i, l := range labels {
					ls[i] = string(l)
				}
				if len(ls) > 0 {
					fmt.Println("orphaned:", strings.Join(ls, ", "))
				}
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
