package config

type LayoutConfig struct {
	BorderWidth      int
	HeaderFooterSize int
	MinPaneWidth     int
	DefaultPaneWidth int
}

type ColorsConfig struct {
	PaneBorder         string
	FocusedPaneBorder  string
	SelectedFileBG     string
	SelectedFileFG     string
	DirFG              string
	FileFG             string
	CmdBoxFG           string
	CmdBoxBorder       string
	AlertNormalFG      string
	AlertNormalBorder  string
	AlertInfoFG        string
	AlertInfoBorder    string
	AlertWarningFG     string
	AlertWarningBorder string
	AlertErrorFG       string
	AlertErrorBorder   string
	ErrorFG            string
	ErrorBG            string
	ErrorBorder        string
}

var Layout = LayoutConfig{
	BorderWidth:      2,
	HeaderFooterSize: 4,
	MinPaneWidth:     20,
	DefaultPaneWidth: 40,
}

var Colors = ColorsConfig{
	PaneBorder:         "62",
	FocusedPaneBorder:  "201",
	SelectedFileBG:     "62",
	SelectedFileFG:     "230",
	DirFG:              "39",
	FileFG:             "252",
	CmdBoxFG:           "230",
	CmdBoxBorder:       "230",
	AlertNormalFG:      "123",
	AlertNormalBorder:  "123",
	AlertInfoFG:        "33",
	AlertInfoBorder:    "33",
	AlertErrorFG:       "9",
	AlertErrorBorder:   "9",
	AlertWarningFG:     "220",
	AlertWarningBorder: "220",
}
