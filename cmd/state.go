package cmd

import (
	"fmt"

	adapter_process "github.com/4sp1/neomux/internal/adapter/os/process"
	adapter_state "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newStateCmd(state adapter_state.Adapter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "manage neomux state",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newStateCleanCmd(state))
	return cmd
}

func newStateCleanCmd(state adapter_state.Adapter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "remove orphan state entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			proc, err := adapter_process.New()
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
