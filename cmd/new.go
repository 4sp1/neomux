package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newNewCmd() (*cobra.Command, error) {
	var rangeStart *int // port start range
	var attach *bool
	var cd *string
	var label string
	cmd := &cobra.Command{
		Use:   "new [LABEL]",
		Short: "creates new nvim server in current directory",
		Args: func(cmd *cobra.Command, args []string) error {
			// if label is not povided read cd and extract base
			if err := cobra.ExactArgs(1)(cmd, args); err != nil {
				label = path.Base(strings.TrimSpace(*cd))
				if label == "" {
					return fmt.Errorf("must provide either label or --cd")
				}
				return nil
			}
			label = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}

			a, err := app.New(nil, state, app.WithMinPort(*rangeStart))
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if err := a.Serve(label, *cd, app.ServeWithAttach(*attach)); err != nil {
				return fmt.Errorf("app: serve: %w", err)
			}

			return nil
		},
	}

	rangeStart = cmd.Flags().Int("range-start", 10010, "minimal port")

	attach = cmd.Flags().Bool("attach", true, "attach to new neovide if true")

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get wd: %w", err)
	}
	cd = cmd.Flags().String("cd", wd, "nvim server root directory")

	return cmd, nil
}
