package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	procs "github.com/4sp1/neomux/internal/adapter/os/process"
	adapter_state "github.com/4sp1/neomux/internal/adapter/sqlite/state"
	"github.com/4sp1/neomux/internal/app"
	"github.com/spf13/cobra"
)

func newListCmd(state adapter_state.Adapter) *cobra.Command {
	var clean, workdir, pid, port, attachedAt *bool
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all nvim servers' labels",
		RunE: func(cmd *cobra.Command, args []string) error {
			proc, err := procs.New()
			if err != nil {
				return fmt.Errorf("procs adapter: new: %w", err)
			}

			app, err := app.New(proc, state)
			if err != nil {
				return fmt.Errorf("app: new: %w", err)
			}

			if *clean {
				labels, err := app.StateClean()
				if err != nil {
					return fmt.Errorf("app: state clean: %w", err)
				}
				{
					ls := make([]string, len(labels))
					for i, l := range labels {
						ls[i] = string(l)
					}
					if len(ls) > 0 {
						fmt.Println("orphaned:", strings.Join(ls, ", "))
					}
				}
			}

			servers, err := state.ListServers(context.Background())
			if err != nil {
				return fmt.Errorf("list servers: %w", err)
			}

			for _, s := range servers {
				info := []string{}
				for _, i := range []struct {
					show  bool
					label string
					value string
				}{
					{*workdir, "workdir", s.Workdir},
					{*pid, "pid", strconv.Itoa(s.PID)},
					{*port, "port", strconv.Itoa(s.Port)},
					{*attachedAt, "attachedAt", s.AttachedAt.Format(time.RFC3339)},
				} {
					var b strings.Builder
					if i.show {
						b.WriteString(i.label)
						b.WriteRune('=')
						b.WriteString(i.value)
					}
					info = append(info, b.String())
				}
				if *workdir || *pid || *port {
					fmt.Printf("label=%s %s\n", s.Label, strings.Join(info, " "))
				} else {
					fmt.Println(s.Label)
				}
			}

			return nil
		},
	}
	clean = cmd.Flags().Bool("clean", false, "clean orphaned sessions from neomux state")
	workdir = cmd.Flags().Bool("wordir", false, "show work directories")
	pid = cmd.Flags().Bool("pid", false, "show nvim headless pids")
	port = cmd.Flags().Bool("port", false, "show nvim headless ports")
	attachedAt = cmd.Flags().Bool("attached-at", false, "show last time neovide was attached to the server")
	return cmd
}
