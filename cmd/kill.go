package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newKillCmd(state adapter.Adapter) *cobra.Command {
	var label *string
	cmd := &cobra.Command{
		Use:   "kill [LABEL]",
		Short: "kill nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				*label = args[0]
			}
			var err error
			if len(*label) == 0 {
				*label, err = fzfRun(cmd.Context(), state)
				if err != nil {
					return fmt.Errorf("fzf: %w", err)
				}
			}

			s, err := state.GetServer(context.Background(), *label)
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
			if err := state.DeleteLabel(context.Background(), *label); err != nil {
				return fmt.Errorf("state: delete label: %w", err)
			}
			return nil
		},
	}
	label = cmd.Flags().String("label", "", "session name")
	return cmd
}
