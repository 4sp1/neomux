package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/4sp1/neomux/internal/domain/server"
	"github.com/4sp1/neomux/internal/domain/workspace"
	"github.com/4sp1/neomux/internal/repo"
	_ "github.com/glebarez/go-sqlite"
)

type adapter struct {
	db *sql.DB
}

func New(path string) (repo.Server, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	return &adapter{
		db: db,
	}, nil
}

func NewWorkspace(path string) (repo.Workspace, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	return &adapter{
		db: db,
	}, nil
}

func (a adapter) Close() error {
	return a.db.Close()
}

var ErrNoLabel = errors.New("label does not exist")

func (a adapter) DeleteLabel(ctx context.Context, label string) error {
	res, err := a.db.ExecContext(ctx, "DELETE FROM neovim_servers WHERE label=?", label)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("delete: %w: %s", ErrNoLabel, label)
	}
	return nil
}

func (a adapter) CreateServer(ctx context.Context, server server.Description) error {
	_, err := a.db.ExecContext(ctx,
		"INSERT INTO neovim_servers (port, pid, label, workdir) VALUES(?, ?, ?, ?)",
		server.Port, server.PID, server.Label, server.Workdir)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func (a adapter) UpdateServerAddr(ctx context.Context, label string, port int, pid int) error {
	_, err := a.db.ExecContext(ctx, "UPDATE neovim_servers SET port = ?, pid = ? WHERE label = ?", port, pid, label)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (a adapter) GetServer(ctx context.Context, label string) (server.Description, error) {
	row := a.db.QueryRow("SELECT port, pid, workdir FROM neovim_servers WHERE label=?", label)
	var s server.Description
	s.Label = label
	if err := row.Scan(&s.Port, &s.PID, &s.Workdir); err != nil {
		return server.Description{}, fmt.Errorf("select: scan: %w", err)
	}
	return s, nil
}

// MaxPort returns (0, nil) if there is no entries in the neovim_servers table.
func (a adapter) MaxPort(ctx context.Context) (int, error) {
	row := a.db.QueryRow("SELECT MAX(port) FROM neovim_servers")
	var maxPort sql.NullInt64
	if err := row.Scan(&maxPort); err != nil {
		return 0, fmt.Errorf("select: scan: %w", err)
	}
	if !maxPort.Valid {
		return 0, nil
	}
	return int(maxPort.Int64), nil
}

func (a adapter) ListServers(ctx context.Context) ([]server.Description, error) {
	rows, err := a.db.Query("SELECT port, pid, label, workdir FROM neovim_servers")
	if err != nil {
		return nil, fmt.Errorf("query: select: %w", err)
	}
	servers := []server.Description{}
	for rows.Next() {
		var s server.Description
		if err := rows.Scan(&s.Port, &s.PID, &s.Label, &s.Workdir); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (a adapter) ListWorkspaces(ctx context.Context) ([]workspace.Description, error) {
	rows, err := a.db.Query("SELECT label, directory FROM workspaces")
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	var workspaces []workspace.Description
	for rows.Next() {
		var w workspace.Description
		if err := rows.Scan(&w.Label, &w.Directory); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		workspaces = append(workspaces, w)
	}
	return workspaces, nil
}

func (a adapter) CreateWorkspace(ctx context.Context, description workspace.Description) error {
	_, err := a.db.Exec("INSERT INTO workspaces (label, directory) VALUES (?, ?)", description.Label, description.Directory)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil

}

func (a adapter) GetWorkspace(ctx context.Context, label string) (*workspace.Description, error) {
	row := a.db.QueryRow("SELECT directory FROM workspaces WHERE label=?", label)
	var directory string
	if err := row.Scan(&directory); err != nil {
		return nil, fmt.Errorf("scan: %w", err)
	}
	return &workspace.Description{
		Label:     label,
		Directory: directory,
	}, nil
}

func (a adapter) DeleteWorkspace(ctx context.Context, label string) error {
	_, err := a.db.ExecContext(ctx, "DELETE FROM workspaces WHERE label = ?", label)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}
