package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type NvimServer struct {
	Port       int
	PID        int
	Label      string
	Workdir    string
	AttachedAt *time.Time
}

type adapter struct {
	db *sql.DB
}

type Adapter interface {
	AttachServer(ctx context.Context, label string, t time.Time) error
	Close() error
	CreateServer(ctx context.Context, server NvimServer) error
	DeleteLabel(ctx context.Context, label string) error
	GetServer(ctx context.Context, label string) (NvimServer, error)
	ListServers(ctx context.Context) ([]NvimServer, error)
	MaxPort(ctx context.Context) (int, error)
	UpdateServerAddr(ctx context.Context, label string, port int, pid int) error
}

func New(path string) (Adapter, error) {
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

func (a adapter) CreateServer(ctx context.Context, server NvimServer) error {
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

func (a adapter) GetServer(ctx context.Context, label string) (NvimServer, error) {
	row := a.db.QueryRow("SELECT port, pid, workdir FROM neovim_servers WHERE label=?", label)
	var s NvimServer
	s.Label = label
	if err := row.Scan(&s.Port, &s.PID, &s.Workdir); err != nil {
		return NvimServer{}, fmt.Errorf("select: scan: %w", err)
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

func (a adapter) ListServers(ctx context.Context) ([]NvimServer, error) {
	rows, err := a.db.Query("SELECT port, pid, label, workdir, attachedAt FROM neovim_servers")
	if err != nil {
		return nil, fmt.Errorf("query: select: %w", err)
	}
	servers := []NvimServer{}
	for rows.Next() {
		var s NvimServer
		var t sql.NullTime
		if err := rows.Scan(&s.Port, &s.PID, &s.Label, &s.Workdir, &t); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		if t.Valid {
			s.AttachedAt = new(time.Time)
			*s.AttachedAt = t.Time
		}
		servers = append(servers, s)
	}
	return servers, nil
}

func (a adapter) AttachServer(ctx context.Context, label string, t time.Time) error {
	_, err := a.db.Exec("UPDATE neovim_servers SET attachedAt = ? WHERE label = ?", t, label)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
