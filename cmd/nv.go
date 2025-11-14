package cmd

import (
	"fmt"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newNvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nv [LABEL]",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
		Short:   "attach neovide to nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}

			app, err := app.New(nil, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if err := app.Attach(args[0]); err != nil {
				return fmt.Errorf("app: attach: %w", err)
			}

			return nil
		},
	}
	return cmd
}
