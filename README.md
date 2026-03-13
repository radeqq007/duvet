<div align="center">
  <img src=".github/assets/logo.png" width="250" alt="Duvet logo">
</div>

# Duvet

Duvet is a terminal based file explorer with vim motions and commands.


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

- `h`, `j`, `k`, `l` - move around
- `<number><motion>` - repeat a motion (e.g. `3j` - go 3 lines down)
- `dd` - delete selected files
- `yy` - yank selected files
- `p` - pasted yanked files
- `gg` - go to the first line
- `G` - go to the last line
- `<line number>G` - go to a specific line
- `Enter` - open the directory or a file
- `Space` - toggle selection
- `Tab` - switch pane focus
- `:` - open the command prompt

## Commands

- `q` - quit
- `quit` - quit
- `rename` - rename current file
- `delete` - delete selected files
- `yank` - yank selected files
- `paste` - paste yanked files
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

## Config


The configuration file is located in the system's default config directory:

- Linux: `~/.config/duvet/config.toml`
- macOS: `~/Library/Application Support/duvet/config.toml`
- Windows: `%APPDATA%\duvet\config.toml`

The default options are as follows:
```toml
default_editor = "vim"
preview_theme = "dracula"

[colors]
pane_border = "159"
focused_pane_border = "153"

selected_file_bg = "62"
selected_file_fg = "230"

dir_fg = "39"
file_fg = "252"

cmd_box_fg = "159"
cmd_box_border = "159"

alert_normal_fg = "123"
alert_normal_border = "123"

alert_info_fg = "33"
alert_info_border = "33"

alert_warning_fg = "220"
alert_warning_border = "220"

alert_error_fg = "9"
alert_error_border = "9"
```

## Demo

![Demo](.github/assets/demo.gif)
