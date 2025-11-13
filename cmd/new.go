package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newNewCmd() (*cobra.Command, error) {
	var rangeStart *int // port start range
	var cd *string
	cmd := &cobra.Command{
		Use:   "new [LABEL]",
		Short: "creates new nvim server in current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("cd", *cd)
			state, err := newState()
			if err != nil {
				return err
			}
			maxPort, err := state.MaxPort(context.Background())
			if err != nil {
				return fmt.Errorf("maxport: %w", err)
			}
			newPort := maxPort + 1
			if newPort < *rangeStart {
				newPort = *rangeStart
			}
			{
				if err := os.Chdir(*cd); err != nil {
					return fmt.Errorf("chdir: %w", err)
				}
				cmd := exec.Command("nvim", "--headless", "--listen", fmt.Sprintf("localhost:%d", newPort))
				cmd.Stdout = nil
				cmd.Stderr = nil
				cmd.Stdin = nil
				if err := cmd.Start(); err != nil {
					return fmt.Errorf("nvim headless start: %w", err)
				}
				label := args[0]
				fmt.Println("new nvim server created with label", label)
				if err := state.CreateServer(context.Background(), adapter.NvimServer{
					PID:   cmd.Process.Pid,
					Label: label,
					Port:  newPort,
				}); err != nil {
					return fmt.Errorf("sqlite: create server: %w", err)
				}
			}
			return nil
		},
	}
	rangeStart = cmd.Flags().Int("range-start", 10010, "minimal port")
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get wd: %w", err)
	}
	cd = cmd.Flags().String("cd", wd, "nvim server root directory")
	return cmd, nil
}
