<div align="center">
  <img src=".github/assets/logo.png" width="250" alt="Duvet logo">
</div>

# What is Duvet?

Duvet is a fast, keyboard-driven file explorer that lives entirely in your terminal. If you use Vim, you already know how to navigate it. Browse directories, preview files with syntax highlighting, manage bookmarks, run shell commands, and see git status — all without leaving the terminal.

## Installation

Requires Go 1.25+ and a **[Nerd Font](https://www.nerdfonts.com/)** for icons.

```sh
go install github.com/radeqq007/duvet/cmd/duvet@latest
```

You run it with:
```sh
duvet
```

## Navigation

| Key | Action |
|-----|--------|
| `h` / `←` | Go to parent directory |
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `l` / `→` / `Enter` | Enter directory or open file |
| `<n><motion>` | Repeat motion `n` times (e.g. `5j`, `3k`) |
| `gg` | Jump to first item |
| `G` | Jump to last item |
| `<n>G` | Jump to line `n` |
| `Space` | Toggle selection on current file |
| `Tab` | Switch pane focus |
| `yy` | Yank selected / current file |
| `p` | Paste yanked files |
| `dd` | Delete selected / current file |
| `Esc` | Clear pending input |
| `:` | Open command prompt |

## Commands

Enter command mode with `:`. Shell commands use the `!` prefix.

### File Operations

| Command | Description |
|---------|-------------|
| `touch <name>` | Create a new file |
| `mkdir <name>` | Create a new directory |
| `rename <name>` | Rename the current file |
| `delete` | Delete selected files |
| `yank` | Yank selected files |
| `paste` | Paste yanked files |
| `cd [path]` | Change directory |

### Selection

| Command | Description |
|---------|-------------|
| `select <pattern>` | Select files matching a glob (e.g. `select *.go`) |
| `deselect <pattern>` | Deselect files matching a glob |

### Search & Bookmarks
 
| Command | Description |
|---------|-------------|
| `find <text>` | Fuzzy-find a file in the current directory and jump to it |
| `bm save <name>` | Save current path as a named bookmark |
| `bm load <name>` | Jump to a saved bookmark |
| `bm list` | List all bookmarks |
| `bm delete <name>` | Remove a bookmark | 

### Shell & Alerts
 
| Command | Description |
|---------|-------------|
| `!<cmd> [args]` | Run a shell command; output shown in an alert box |
| `alert <text>` | Show a normal alert |
| `alert <type> <text>` | Show a typed alert (`normal`, `info`, `warning`, `error`) |
 
### App
 
| Command | Description |
|---------|-------------|
| `q` / `quit` | Quit |


## Git integration

Duvet shows git status inline next to file names when you're inside a git repository. The status codes use the standard two-character format from `git status --porcelain`.
 
 
The current branch name is also shown in the status bar at the bottom.

## File Preview
 
Selecting a file automatically loads a syntax-highlighted preview in the right pane.
 
The preview theme is configurable (default: `dracula`). Any [Chroma theme](https://xyproto.github.io/splash/docs/) works.

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
