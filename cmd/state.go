package cmd

import (
	"fmt"

	adapter "github.com/4sp1/neomux/internal/adapter/os/process"
	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "manage neomux state",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newStateCleanCmd())
	return cmd
}

func newStateCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove orphan state entries",
		RunE: func(cmd *cobra.Command, args []string) error {

			state, err := newState()
			if err != nil {
				return err
			}

			proc, err := adapter.New()
			if err != nil {
				return err
			}

			app, err := app.New(proc, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			labels, err := app.StateClean()
			if err != nil {
				return fmt.Errorf("app: state clean: %w", err)
			}

			for _, label := range labels {
				fmt.Println(label)
			}

			return nil
		},
	}
	return cmd
}
