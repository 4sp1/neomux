package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func newNvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nv [LABEL]",
		Aliases: []string{"a"},
		Args:    cobra.ExactArgs(1),
		Short:   "spawns new neovide to nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}
			label := args[0]
			s, err := state.GetServer(context.Background(), label)
			if err != nil {
				return fmt.Errorf("get neovim server: %w", err)
			}
			{
				cmd := exec.Command("neovide", "--frame=transparent", fmt.Sprintf("--server=localhost:%d", s.Port))
				if err := cmd.Start(); err != nil {
					return fmt.Errorf("run neovide: %w", err)
				}
			}
			return nil
		},
	}
	return cmd
}
