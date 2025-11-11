package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newKillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill [LABEL]",
		Short: "kill nvim server",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := statePath()
			if err != nil {
				return fmt.Errorf("state path: %w", err)
			}
			state, err := adapter.New(path)
			if err != nil {
				return fmt.Errorf("new sqlite state adapter: %w", err)
			}
			label := args[0]
			s, err := state.GetServer(context.Background(), label)
			if err != nil {
				return fmt.Errorf("state: get server: %w", err)
			}
			{
				cmd := exec.Command("kill", fmt.Sprintf("%d", s.PID))
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("kill command: %w", err)
				}
			}
			if err := state.DeleteLabel(context.Background(), label); err != nil {
				return fmt.Errorf("state: delete label: %w", err)
			}
			return nil
		},
	}
	return cmd
}
