package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type fileConfig struct {
	Colors        ColorsConfig `toml:"colors"`
	DefaultEditor string       `toml:"default_editor"`
	PreviewTheme  string       `toml:"preview_theme"`
}

type LayoutConfig struct {
	BorderWidth      int
	HeaderFooterSize int
	MinPaneWidth     int
	DefaultPaneWidth int
}

type ColorsConfig struct {
	PaneBorder         string `toml:"pane_border"`
	FocusedPaneBorder  string `toml:"focused_pane_border"`
	SelectedFileBG     string `toml:"selected_file_bg"`
	SelectedFileFG     string `toml:"selected_file_fg"`
	DirFG              string `toml:"dir_fg"`
	FileFG             string `toml:"file_fg"`
	CmdBoxFG           string `toml:"cmd_box_fg"`
	CmdBoxBorder       string `toml:"cmd_box_border"`
	AlertNormalFG      string `toml:"alert_normal_fg"`
	AlertNormalBorder  string `toml:"alert_normal_border"`
	AlertInfoFG        string `toml:"alert_info_fg"`
	AlertInfoBorder    string `toml:"alert_info_border"`
	AlertWarningFG     string `toml:"alert_warning_fg"`
	AlertWarningBorder string `toml:"alert_warning_border"`
	AlertErrorFG       string `toml:"alert_error_fg"`
	AlertErrorBorder   string `toml:"alert_error_border"`
	ErrorFG            string `toml:"error_fg"`
	ErrorBG            string `toml:"error_bg"`
	ErrorBorder        string `toml:"error_border"`
}

var Layout = LayoutConfig{
	BorderWidth:      2,
	HeaderFooterSize: 4,
	MinPaneWidth:     20,
	DefaultPaneWidth: 40,
}

var Colors = ColorsConfig{
	PaneBorder:         "159",
	FocusedPaneBorder:  "153",
	SelectedFileBG:     "62",
	SelectedFileFG:     "230",
	DirFG:              "39",
	FileFG:             "252",
	CmdBoxFG:           "159",
	CmdBoxBorder:       "159",
	AlertNormalFG:      "123",
	AlertNormalBorder:  "123",
	AlertInfoFG:        "33",
	AlertInfoBorder:    "33",
	AlertErrorFG:       "9",
	AlertErrorBorder:   "9",
	AlertWarningFG:     "220",
	AlertWarningBorder: "220",
}

var (
	DefaultEditor = "vim"
	PreviewTheme  = "dracula"
)

func Load() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "duvet", "config.toml")
	if _, err := os.Stat(configPath); err != nil && os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	var fc fileConfig
	fc.Colors = Colors
	fc.DefaultEditor = DefaultEditor
	fc.PreviewTheme = PreviewTheme

	if err := toml.Unmarshal(data, &fc); err != nil {
		return err
	}

	Colors = fc.Colors
	DefaultEditor = fc.DefaultEditor
	PreviewTheme = fc.PreviewTheme

	return nil
}
