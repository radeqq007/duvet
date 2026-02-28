<div align="center">
  <img src=".github/assets/logo.png" width="250" alt="Duvet logo">
</div>

# Duvet

Duvet is a terminal based file explorer with vim inspired motions and commands.


## Installation

You can install duvet by running:
```sh
go install github.com/radeqq007/duvet/cmd/duvet@latest
```

and you can run it with:
```sh
duvet
```

In order for icons to render properly you also need a **[Nerd Font](https://www.nerdfonts.com/)** installed.

## Navigation

- `h`, `j`, `k`, `l` or the arrow keys - move around
- `Enter` - open the directory or a file
- `Space` - toggle selection
- `:` - open the command prompt

## Commands

- `q` - quit
- `quit` - quit
- `rename` - rename current file
- `delete` - delete selected files
- `touch <file name>` - create a new file
- `mkdir <directory name>` - create a new directory
- `cd` - change directory
- `bm save <name>` - save the current path as a bookmark
- `bm delete <name>` - remove a bookmark
- `bm load <name>` - load a bookmark
- `bm list` - list all saved bookmarks
- `find <text>` - fuzzy match a file from current directory and jump to it
- `alert <type> <text>` - open the alert box (possible types: `normal`, `info`, `warning`, `error`)
- `alert <text>` - open the alert box, with the default type `normal`
- `select <pattern>` - select files that match the pattern (e.g. `:select *.py`)
- `deselect <pattern>` - deselect files that match the pattern (e.g. `:deselect *.rs`)
- `!<command> <args>` - executes a shell command

## Demo

![Demo](.github/assets/demo.gif)
