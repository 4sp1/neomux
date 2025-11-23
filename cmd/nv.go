package cmd

import (
	"fmt"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newNvCmd() *cobra.Command {
	var label *string
	cmd := &cobra.Command{
		Use:     "attach [LABEL]",
		Aliases: []string{"a", "nv"},
		Short:   "attach neovide to nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				*label = args[0]
			}

			if *label == "" {
				*label, err = fzfRun(cmd.Context(), state)
				if err != nil {
					return fmt.Errorf("fzf: %w", err)
				}
			}

			app, err := app.New(nil, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if err := app.Attach(*label); err != nil {
				return fmt.Errorf("app: attach: %w", err)
			}

			return nil
		},
	}
	label = cmd.Flags().String("label", "", "session name")
	return cmd
}
