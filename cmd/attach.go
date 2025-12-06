package cmd

import (
	"fmt"

	proc_adapter "github.com/4sp1/neomux/internal/adapter/os/process"
	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newNvCmd() *cobra.Command {
	var label *string
	var noReset *bool
	cmd := &cobra.Command{
		Use:     "attach [LABEL]",
		Aliases: []string{"a", "nv"},
		Short:   "attach neovide to nvim server",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := newState()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				*label = args[0]
			}

			if *label == "" {
				*label, err = fzfRun(cmd.Context(), state)
				if err != nil {
					return fmt.Errorf("fzf: %w", err)
				}
			}

			var proc proc_adapter.Adapter
			if !*noReset {
				proc, err = proc_adapter.New()
				if err != nil {
					return fmt.Errorf("new proc adapter: %w", err)
				}
			}

			app, err := app.New(proc, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if *noReset {
				if err := app.Attach(*label); err != nil {
					return fmt.Errorf("app: attach: %w", err)
				}
				return nil
			}

			if err := app.AttachOrRestore(*label); err != nil {
				return fmt.Errorf("app: attach or restore: %w", err)
			}

			return nil
		},
	}
	label = cmd.Flags().String("label", "", "session name")
	noReset = cmd.Flags().Bool("no-reset", false, "no new instance of nvim server created")
	return cmd
}
