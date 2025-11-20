package cmd

import (
	"fmt"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newDuplicateCmd() *cobra.Command {
	var attach *bool
	cmd := &cobra.Command{
		Use:     "duplicate [LABEL]",
		Aliases: []string{"clone", "dup"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := newState()
			if err != nil {
				return fmt.Errorf("new state: %w", err)
			}

			a, err := app.New(nil, s, app.WithDebug())
			if err != nil {
				return fmt.Errorf("new app: %w", err)
			}

			label, err := a.Duplicate(args[0], app.ServeWithAttach(*attach))
			if err != nil {
				return fmt.Errorf("app: duplicate: %w", err)
			}

			fmt.Println(label)
			return nil
		},
	}
	attach = cmd.Flags().Bool("attach", true, "attach to new neovide if true")
	return cmd
}
