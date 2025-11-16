# Neomux

Neovim/Neovide multiplexer.

Neovim: A highly configurable text editor built on Vim with modern
features and extensibility. Neovim's headless mode runs without a UI,
enabling automation, scripting, and server-like operations through its
API and RPC interface.

Neovide: A graphical user interface for Neovim written in
Rust. It provides visual enhancements over the terminal UI while
maintaining functional compatibility, offering a modern cross-platform
editing experience.

Neomux: Manages headless Neovim instances and spawns Neovide.

## Requirements

- `nc`: netcat openbsd variant
- `nvim`
- `neovide`

## Commands

```
neomux
    new [label]     -- creates new nvim server in current directory
    nv [label]      -- spawns new neovide to nvim server
    kill [label]    -- kill nvim server
    list            -- list all nvim servers' labels
```

## Servers

We will use port ranges starting from port 10000. 
Currently port collision is not supported.

## State

`state.db`

```sql
CREATE TABLE port_table (
    port INTEGER NOT NULL,
    label TEXT NOT NULL,
    nvim_pid INTEGER NOT NULL
);
```

## TODO

### First:

- kill does not check for return status of /bin/kill
- better error handling 
  - prettier errors for "no rows"
  - "unique constraint"
- handle case where nvim server killed by another process "No such process"
  (maybe refresh command, e.g. system reboots)
- update README.md with installation instructions and workflow demo

### Then:

- ssh tunneling (local port forwarding) and add command for next port
- state clean command support for other platforms (list processes)
  see [mitchellh/go-ps](
  https://github.com/mitchellh/go-ps
  )
