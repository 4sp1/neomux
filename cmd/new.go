package cmd

import (
	"fmt"
	"os"

	"github.com/4sp1/neomux/internal/app"
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
			state, err := newState()
			if err != nil {
				return err
			}

			app, err := app.New(nil, state, app.WithMinPort(*rangeStart))
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if err := app.Serve(args[0], *cd); err != nil {
				return fmt.Errorf("app: serve: %w", err)
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
