package icons

type Icon struct {
	Icon  string
	Color string
}

//var icons map[string]string = map[string]string{
//	".cpp":       "\ue61d",
//	".cs":        "\ue7b2",
//	".html":      "\ue60e",
//	".htm":       "\ue60e",
//	".css":       "\ue6b8",
//	".gitignore": "\ue702",
//	".py":        "\ue73c",
//	".php":       "\ue73d",
//	".java":      "\ue738",
//	".rb":        "\ue605",
//	".swift":     "\ue699",
//	".kt":        "\ue634",
//	".yaml":      "\ue8eb",
//	".yml":       "\ue8eb",
//	".sh":        "\ue691",
//	".md":        "\ueb1d",
//	".jpg":       "\uf03e",
//	".jpeg":      "\uf03e",
//	".png":       "\uf03e",
//	".webp":      "\uf03e",
//	".gif":       "\uf0d78",
//	".svg":       "\uf0721",
//	".pdf":       "\uf1c1",
//	".mp3":       "\ue638",
//	".wav":       "\ue638",s
//}

var icons map[string]Icon = map[string]Icon{
	".cpp": {
		Color: "20",
		Icon:  "\ue61d",
	},
	".h": {
		Color: "20",
		Icon:  "\ue61e",
	},
	".c": {
		Color: "20",
		Icon:  "\ue61e",
	},
	".rs": {
		Color: "202",
		Icon:  "\ue7a8",
	},
	".json": {
		Color: "220",
		Icon:  "\ueb0f",
	},
	".ts": {
		Color: "27",
		Icon:  "\ue69d",
	},
	".js": {
		Color: "220",
		Icon:  "\uf2ee",
	},
	".mod": {
		Color: "39",
		Icon:  "\ue627",
	},
	".sum": {
		Color: "39",
		Icon:  "\ue627",
	},
}

func GetIcon(fileExtension string) Icon {
	if icon, ok := icons[fileExtension]; ok {
		return icon
		// return lipgloss.NewStyle().Foreground(lipgloss.Color(icon.Color)).Render(icon.Icon)
	}

	return Icon{}
	// return "\uf15b"
}
