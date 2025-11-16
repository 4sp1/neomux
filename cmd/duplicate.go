package cmd

import (
	"fmt"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newDuplicateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "duplicate [LABEL]",
		Aliases: []string{"clone", "dup"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := newState()
			if err != nil {
				return fmt.Errorf("new state: %w", err)
			}

			app, err := app.New(nil, s, app.WithDebug())
			if err != nil {
				return fmt.Errorf("new app: %w", err)
			}

			label, err := app.Duplicate(args[0])
			if err != nil {
				return fmt.Errorf("app: duplicate: %w", err)
			}

			fmt.Println(label)
			return nil
		},
	}
	return cmd
}
