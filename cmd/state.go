package cmd

import (
	"context"
	"fmt"

	adapter "github.com/4sp1/neomux/internal/adapter/os/process"
	adapter_state "github.com/4sp1/neomux/internal/adapter/sqlite/state"
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
			servers, err := state.ListServers(context.Background())
			if err != nil {
				return fmt.Errorf("list servers: %w", err)
			}
			processes, err := adapter.New()
			if err != nil {
				return err
			}
			procs, err := processes.List()
			if err != nil {
				return fmt.Errorf("processes list: %w", err)
			}
			pmap := map[int]string{}
			for _, p := range procs {
				pmap[p.PID] = p.Binary
			}
			for _, s := range servers {
				if _, ok := pmap[s.PID]; !ok {
					if err := deleteLabel(state, s.Label); err != nil {
						return err
					}
				}
				if b, ok := pmap[s.PID]; ok && b != "nvim" {
					if err := deleteLabel(state, s.Label); err != nil {
						return err
					}
				}
			}
			return nil
		},
	}
	return cmd
}

func deleteLabel(state adapter_state.Adapter, label string) error {
	if err := state.DeleteLabel(context.Background(), label); err != nil {
		return fmt.Errorf("delete label %q: %w", label, err)
	}
	fmt.Println(label)
	return nil
}
