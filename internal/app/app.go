package app

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/4sp1/neomux/internal/domain/server"
	"github.com/4sp1/neomux/internal/domain/workspace"
	"github.com/4sp1/neomux/internal/repo"
)

func New(p repo.Proc, s repo.Server, w repo.Workspace, opts ...Option) (App, error) {
	var c Config
	c.minPort = 10000
	c = OptionTermUI()(c)
	for _, opt := range opts {
		c = opt(c)
	}
	return &app{
		proc:      p,
		server:    s,
		workspace: w,
		conf:      c,
	}, nil
}

type Label string

type App interface {
	Serve(label, workdir string, opts ...ServeOption) error
	Attach(label string) error
	AttachOrRestore(label string) error
	StateClean() ([]Label, error)
	Duplicate(label string, opts ...ServeOption) (string, error)
	Shell(ctx context.Context, description workspace.Description) error
	ListWorkspaces(ctx context.Context) error
}

type app struct {
	proc      repo.Proc
	server    repo.Server
	workspace repo.Workspace
	conf      Config
}

type Config struct {
	minPort int
	debug   bool

	cmd  func(port int) *exec.Cmd
	fork bool
}

type Option func(Config) Config

func OptionDebug() Option {
	return func(c Config) Config {
		c.debug = true
		return c
	}
}

func OptionStartPort(port int) Option {
	return func(c Config) Config {
		c.minPort = port
		return c
	}
}

func OptionNeovideUI(columns, lines int) Option {
	if columns == 0 {
		columns = 160
	}
	if lines == 0 {
		lines = 120
	}
	return OptionUI(true, "neovide",
		"--frame-transparent",
		fmt.Sprintf("--grid=%dx%d", lines, columns),
		"--server=localhost:%d")
}

func OptionTermUI() Option {
	return OptionUI(false, "nvim",
		"--server", "localhost:%d",
		"--remote-ui")
}

func OptionUI(fork bool, command string, args ...string) Option {
	return func(c Config) Config {
		c.cmd = func(port int) *exec.Cmd {
			for i, arg := range args {
				if strings.Contains(arg, "%d") {
					args[i] = fmt.Sprintf(arg, port)
				}
			}
			return exec.Command(command, args...)
		}
		c.fork = fork
		return c
	}
}

type ServeOption func(ServeConfig) ServeConfig
type ServeConfig struct {
	attach  bool
	restore bool
}

func ServeWithAttach(attach bool) ServeOption {
	return func(sc ServeConfig) ServeConfig {
		sc.attach = attach
		return sc
	}
}

func ServeWithRestore(restore bool) ServeOption {
	return func(sc ServeConfig) ServeConfig {
		sc.restore = restore
		return sc
	}
}

func (a app) ListWorkspaces(ctx context.Context) error {
	workspaces, err := a.workspace.ListWorkspaces(ctx)
	if err != nil {
		return fmt.Errorf("list workspaces: %w", err)
	}
	var maxLen int
	for _, workspace := range workspaces {
		if maxLen < len(workspace.Label) {
			maxLen = len(workspace.Label)
		}
	}
	for _, w := range workspaces {
		var b strings.Builder
		b.WriteString(w.Label)
		for b.Len() < maxLen+2 {
			b.WriteRune(' ')
		}
		b.WriteString(w.Directory)
		fmt.Println(b.String())
	}
	return nil
}

func (a app) getOrCreateWorkspace(ctx context.Context, description *workspace.Description) error {
	lookup, err := a.workspace.GetWorkspace(ctx, description.Label)
	if err != nil && description.Directory == "" {
		return fmt.Errorf("no directory provided: get workspace: %w", err)
	}
	if err != nil && description.Directory != "" {
		if err := a.workspace.CreateWorkspace(ctx, *description); err != nil {
			return fmt.Errorf("create workspace: %w", err)
		}
		return nil
	}
	description.Directory = lookup.Directory
	return nil
}

func (a app) Shell(ctx context.Context, description workspace.Description) error {
	if err := a.getOrCreateWorkspace(ctx, &description); err != nil {
		return fmt.Errorf("shell: %w", err)
	}
	return shell(description.Directory)
}

func shell(dir string) error {
	shell, found := os.LookupEnv("SHELL")
	if !found {
		return errors.New("the required SHELL environment variable could not be located")
	}
	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("chdir %q: %w", dir, err)
	}
	cmd := exec.Command(shell)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (a app) Serve(label, workdir string, opts ...ServeOption) error {
	conf := ServeConfig{}
	for _, opt := range opts {
		conf = opt(conf)
	}

	newPort := a.conf.minPort
	for {
		err := exec.Command("nc", "-z", "localhost", fmt.Sprintf("%d", newPort)).Run()
		if err == nil {
			if a.conf.debug {
				fmt.Println("port", newPort, "in use, trying", newPort+1)
			}
			newPort += 1
			continue
		} else {
			switch err := err.(type) {
			case *exec.ExitError:
			default:
				return fmt.Errorf("exec: netcat: %w", err)
			}
		}
		break
	}

	if !path.IsAbs(workdir) {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("os: getwd: %w", err)
		}
		workdir = path.Join(wd, workdir)
	}

	if err := os.Chdir(workdir); err != nil {
		return fmt.Errorf("chdir: %w", err)
	}

	cmd := exec.Command("nvim", "--headless", "--listen", fmt.Sprintf("localhost:%d", newPort))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec nvim headless localhost:%d: %w", newPort, err)
	}

	if a.conf.debug {
		fmt.Println("new server", "port =", newPort, "pid =", cmd.Process.Pid)
	}

	if conf.restore {
		if err := a.server.UpdateServerAddr(context.Background(), label, newPort, cmd.Process.Pid); err != nil {
			return fmt.Errorf("state: update server port: %w", err)
		}
	} else {
		if err := a.server.CreateServer(context.Background(), server.Description{
			PID:     cmd.Process.Pid,
			Label:   label,
			Port:    newPort,
			Workdir: workdir,
		}); err != nil {
			return fmt.Errorf("state: create server: %w", err)
		}
	}

	if conf.attach {
		if err := a.Attach(label); err != nil {
			return fmt.Errorf("app: attach: %w", err)
		}
	}

	return nil
}

func (a app) AttachOrRestore(label string) error {
	s, err := a.server.GetServer(context.TODO(), label)
	if err != nil {
		return fmt.Errorf("state: get server %q: %w", label, err)
	}

	procs, err := a.proc.List()
	if err != nil {
		return fmt.Errorf("proc: list: %w", err)
	}

	// search in server process in procs and attach if found
	for _, p := range procs {
		if a.conf.debug {
			fmt.Println(p.PID, p.Binary, s.PID, s.Label)
		}
		if p.PID == s.PID && p.Binary == "nvim" {
			return a.Attach(label)
		}
	}

	// if no nvim server was found, start a new server
	// and attach a new neovide session
	return a.Serve(s.Label, s.Workdir, ServeWithAttach(true), ServeWithRestore(true))
}

func (a app) Attach(label string) error {
	s, err := a.server.GetServer(context.TODO(), label)
	if err != nil {
		return fmt.Errorf("state: get server %q: %w", label, err)
	}

	cmd := a.conf.cmd(s.Port)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	mode := cmd.Run
	if a.conf.fork {
		mode = cmd.Start
	}

	if err := mode(); err != nil {
		return fmt.Errorf("exec: command neovide: %w", err)
	}

	return nil
}
func (a app) StateClean() ([]Label, error) {
	servers, err := a.server.ListServers(context.Background())
	if err != nil {
		return nil, fmt.Errorf("state: list servers: %w", err)
	}

	procs, err := a.proc.List()
	if err != nil {
		return nil, fmt.Errorf("proc: list: %w", err)
	}

	pmap := map[int]string{}
	for _, p := range procs {
		pmap[p.PID] = p.Binary
	}

	deleted := make([]Label, 0, len(servers))

	for _, s := range servers {
		if _, ok := pmap[s.PID]; !ok {
			if err := a.deleteLabel(s.Label); err != nil {
				return nil, err
			}
			deleted = append(deleted, Label(s.Label))
		}
		if b, ok := pmap[s.PID]; ok && b != "nvim" {
			if err := a.deleteLabel(s.Label); err != nil {
				return nil, err
			}
			deleted = append(deleted, Label(s.Label))
		}
	}

	return deleted, nil

}

func (a app) deleteLabel(label string) error {
	if err := a.server.DeleteLabel(context.Background(), label); err != nil {
		return fmt.Errorf("state: delete label %q: %w", label, err)
	}
	return nil
}

func (a app) Duplicate(label string, opts ...ServeOption) (string, error) {
	conf := ServeConfig{}
	for _, opt := range opts {
		conf = opt(conf)
	}

	s, err := a.server.GetServer(context.Background(), label)
	if err != nil {
		return "", fmt.Errorf("state: get server: %w", err)
	}
	var newLabel string
	{
		charset := []byte{
			'a', 'b', 'c', 'd', 'e', 'x', 'y', 'z',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0',
			'A', 'B', 'C', 'D', 'E', 'X', 'Y', 'Z'}
		b := make([]byte, 4)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		newLabel = label + "-" + string(b)
	}
	if err := a.Serve(newLabel, s.Workdir, ServeWithAttach(conf.attach)); err != nil {
		return "", fmt.Errorf("serve %q: %w", newLabel, err)
	}
	return newLabel, nil
}
