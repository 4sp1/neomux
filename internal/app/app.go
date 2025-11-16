package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	proc_adapter "github.com/4sp1/neomux/internal/adapter/os/process"
	state_adapter "github.com/4sp1/neomux/internal/adapter/sqlite/state"
)

func New(p proc_adapter.Adapter, s state_adapter.Adapter, opts ...Option) (App, error) {
	var c Config
	c.minPort = 10000
	for _, opt := range opts {
		c = opt(c)
	}
	return &app{
		proc:  p,
		state: s,
		conf:  c,
	}, nil
}

type Label string

type App interface {
	Serve(label, workdir string) error
	Attach(label string) error
	StateClean() ([]Label, error)
	Duplicate(label string) (string, error)
}

type app struct {
	proc  proc_adapter.Adapter
	state state_adapter.Adapter
	conf  Config
}

type Config struct {
	minPort int
}

type Option func(Config) Config

func WithMinPort(port int) Option {
	return func(c Config) Config {
		c.minPort = port
		return c
	}
}

func (a app) Serve(label, workdir string) error {
	maxPort, err := a.state.MaxPort(context.TODO())
	if err != nil {
		return fmt.Errorf("maxport: %w", err)
	}
	newPort := maxPort + 1
	if newPort < a.conf.minPort {
		newPort = a.conf.minPort
	}

	if err := os.Chdir(workdir); err != nil {
		return fmt.Errorf("chdir: %w", err)
	}

	cmd := exec.Command("nvim", "--headless", "--listen", fmt.Sprintf("localhost:%d", newPort))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec nvim headless localhost:%d: %w", newPort, err)
	}

	if err := a.state.CreateServer(context.Background(), state_adapter.NvimServer{
		PID:     cmd.Process.Pid,
		Label:   label,
		Port:    newPort,
		Workdir: workdir,
	}); err != nil {
		return fmt.Errorf("state: create server: %w", err)
	}

	return nil
}

func (a app) Attach(label string) error {
	s, err := a.state.GetServer(context.TODO(), label)
	if err != nil {
		return fmt.Errorf("state: get server %q: %w", label, err)
	}

	cmd := exec.Command("neovide",
		"--frame=transparent", "--grid=120x80",
		fmt.Sprintf("--server=localhost:%d", s.Port))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("exec: command neovide: %w", err)
	}

	return nil
}
func (a app) StateClean() ([]Label, error) {
	servers, err := a.state.ListServers(context.Background())
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
	if err := a.state.DeleteLabel(context.Background(), label); err != nil {
		return fmt.Errorf("state: delete label %q: %w", label, err)
	}
	return nil
}

func (a app) Duplicate(label string) (string, error) {
	s, err := a.state.GetServer(context.Background(), label)
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
	if err := a.Serve(newLabel, s.Workdir); err != nil {
		return "", fmt.Errorf("serve %q: %w", newLabel, err)
	}
	return newLabel, nil
}
