package repo

import (
	"context"

	"github.com/4sp1/neomux/internal/domain/proc"
	"github.com/4sp1/neomux/internal/domain/server"
	"github.com/4sp1/neomux/internal/domain/workspace"
)

type Server interface {
	DeleteLabel(ctx context.Context, label string) error
	CreateServer(ctx context.Context, server server.Description) error
	UpdateServerAddr(ctx context.Context, label string, port int, pid int) error
	GetServer(ctx context.Context, label string) (server.Description, error)
	MaxPort(ctx context.Context) (int, error)
	ListServers(ctx context.Context) ([]server.Description, error)
	Close() error
}

type Workspace interface {
	CreateWorkspace(ctx context.Context, description workspace.Description) error
	GetWorkspace(ctx context.Context, label string) (*workspace.Description, error)
	ListWorkspaces(ctx context.Context) ([]workspace.Description, error)
	DeleteWorkspace(ctx context.Context, label string) error
}

type Proc interface {
	List() ([]proc.Proc, error)
}
