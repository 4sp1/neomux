package cmd

import (
	"fmt"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newDuplicateCmd() *cobra.Command {
	var attach *bool
	var label *string
	cmd := &cobra.Command{
		Use:     "duplicate [LABEL]",
		Aliases: []string{"clone", "dup"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				*label = args[0]
			}

			s, err := newState()
			if err != nil {
				return fmt.Errorf("new state: %w", err)
			}

			if len(*label) == 0 {
				*label, err = fzfRun(cmd.Context(), s)
				if err != nil {
					return fmt.Errorf("fzf: %w", err)
				}
			}

			a, err := app.New(nil, s, app.WithDebug())
			if err != nil {
				return fmt.Errorf("new app: %w", err)
			}

			*label, err = a.Duplicate(*label, app.ServeWithAttach(*attach))
			if err != nil {
				return fmt.Errorf("app: duplicate: %w", err)
			}

			fmt.Println(*label)
			return nil
		},
	}
	attach = cmd.Flags().Bool("attach", true, "attach to new neovide if true")
	label = cmd.Flags().String("label", "", "session name")
	return cmd
}
