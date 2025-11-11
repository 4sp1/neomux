CREATE TABLE neovim_servers (
	port INTEGER NOT NULL,
	pid INTEGER NOT NULL,
	label TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_port_table_label_unique ON neovim_servers (label);
