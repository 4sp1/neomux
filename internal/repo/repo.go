package repo

import (
	"context"

	"github.com/4sp1/neomux/internal/domain/proc"
	"github.com/4sp1/neomux/internal/domain/server"
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

type Proc interface {
	List() ([]proc.Proc, error)
}
