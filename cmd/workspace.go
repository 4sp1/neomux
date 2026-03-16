package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/4sp1/neomux/internal/app"
	"github.com/4sp1/neomux/internal/domain/workspace"
	"github.com/spf13/cobra"
)

func newWorkspaceCommand() (*cobra.Command, error) {
	var directory *string

	state, path, err := newStateWorkspace()
	if err != nil {
		return nil, err
	}

	app, err := app.New(nil, nil, state)
	if err != nil {
		return nil, err
	}

	label := "..."
	if len(os.Args) == 2 {
		label = os.Args[1]
	}

	cmd := &cobra.Command{
		Args:    cobra.ExactArgs(1),
		Use:     "workspace [-d DIR| --directory DIR] LABEL",
		Aliases: []string{"w"},
		Long: fmt.Sprintf(`
The command creates a workspace labeled LABEL.

If that workspace does not already exist, it links the specified directory
(defaulting to the current directory) to it.

A workspace can be created only once; to delete it

neomux workspace delete %s
`, label),
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Shell(context.Background(), workspace.Description{
				Label:     args[0],
				Directory: *directory,
			})
		},
	}
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	directory = cmd.Flags().StringP("directory", "d", dir, "create workspace in directory")

	cmd.AddCommand(newWorkspaceListCommand(app, path))
	cmd.AddCommand(newWorkspaceDeleteCommand(app))

	return cmd, nil
}

func newWorkspaceListCommand(app app.App, statePath string) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.ListWorkspaces(context.Background()); err != nil {
				fmt.Println("state db:", statePath)
				return err
			}
			return nil
		},
	}
}

func newWorkspaceDeleteCommand(app app.App) *cobra.Command {
	return &cobra.Command{
		Use:  "delete LABEL [LABEL...]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, label := range args {
				if err := app.DeleteWorkspace(cmd.Context(), label); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
