package cmd

import (
	"context"
	"fmt"
	"os/exec"

	adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	var rangeStart *int // port start range
	cmd := &cobra.Command{
		Use:   "new [LABEL]",
		Short: "creates new nvim server in current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := statePath()
			if err != nil {
				return fmt.Errorf("state path: %w", err)
			}
			state, err := adapter.New(path)
			if err != nil {
				return fmt.Errorf("new sqlite adapter: %w (path=%q)", err, path)
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
	return cmd
}
