package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Colors        ColorsConfig `toml:"colors"`
	Layout        LayoutConfig `toml:"-"`
	DefaultEditor string       `toml:"default_editor"`
	PreviewTheme  string       `toml:"preview_theme"`
}

type LayoutConfig struct {
	BorderWidth      int
	HeaderFooterSize int
	StatusBarHeight  int
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

var defaultLayout = LayoutConfig{
	BorderWidth:      1,
	HeaderFooterSize: 4,
	StatusBarHeight:  1,
	MinPaneWidth:     20,
	DefaultPaneWidth: 40,
}

var defaultColors = ColorsConfig{
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
	defaultEditor       = "vim"
	defaultPreviewTheme = "dracula"
)

func Get() (*Config, error) {

	var conf Config
	conf.Colors = defaultColors
	conf.Layout = defaultLayout
	conf.DefaultEditor = defaultEditor
	conf.PreviewTheme = defaultPreviewTheme

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "duvet", "config.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &conf, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if err := toml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
