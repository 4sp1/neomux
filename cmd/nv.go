package cmd

import (
	"context"
	"fmt"
	"os/exec"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newNvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nv [LABEL]",
		Args:  cobra.ExactArgs(1),
		Short: "spawns new neovide to nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := statePath()
			if err != nil {
				return fmt.Errorf("state path: %w", err)
			}
			state, err := adapter.New(path)
			if err != nil {
				return fmt.Errorf("sqlite state adapter: %w", err)
			}
			label := args[0]
			s, err := state.GetServer(context.Background(), label)
			if err != nil {
				return fmt.Errorf("get neovim server: %w", err)
			}
			{
				cmd := exec.Command("neovide", fmt.Sprintf("--server=localhost:%d", s.Port))
				if err := cmd.Start(); err != nil {
					return fmt.Errorf("run neovide: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}
