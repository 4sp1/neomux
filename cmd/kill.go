package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newKillCmd() *cobra.Command {
	var label *string
	cmd := &cobra.Command{
		Use:   "kill [LABEL]",
		Short: "kill nvim server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			state, err := newState()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				*label = args[0]
			}

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

			defer func() {
				deleteErr := state.DeleteLabel(context.Background(), *label)
				if deleteErr != nil {
					fmt.Println(*label, "removed from state")
				}
				err = deferErr(err, deleteErr)
			}()

			{
				cmd := exec.Command("kill", fmt.Sprintf("%d", s.PID))
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("kill command: %w", err)
				}
			}
			return nil
		},
	}
	label = cmd.Flags().String("label", "", "session name")
	return cmd
}

func deferErr(err, deferErr error) error {
	if err != nil && deferErr != nil {
		return fmt.Errorf("%w then %w", err, deferErr)
	}
	if deferErr != nil {
		return deferErr
	}
	return err
}
