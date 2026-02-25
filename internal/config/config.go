package config

type LayoutConfig struct {
	BorderWidth      int
	HeaderFooterSize int
	MinPaneWidth     int
	DefaultPaneWidth int
}

type ColorsConfig struct {
	PaneBorder        string
	FocusedPaneBorder string
	SelectedFileBG    string
	SelectedFileFG    string
	DirFG             string
	FileFG            string
	CmdBoxFG          string
	CmdBoxBG          string
	CmdBoxBorder      string
	AlertFG           string
	AlertBG           string
	AlertBorder       string
	ErrorFG           string
	ErrorBG           string
	ErrorBorder       string
}

var Layout = LayoutConfig{
	BorderWidth:      2,
	HeaderFooterSize: 4,
	MinPaneWidth:     20,
	DefaultPaneWidth: 40,
}

var Colors = ColorsConfig{
	PaneBorder:        "62",
	FocusedPaneBorder: "201",
	SelectedFileBG:    "62",
	SelectedFileFG:    "230",
	DirFG:             "39",
	FileFG:            "252",
	CmdBoxFG:          "230",
	CmdBoxBG:          "230",
	CmdBoxBorder:      "230",
	AlertFG:           "9",
	AlertBG:           "9",
	AlertBorder:       "9",
}
