package icons

var icons map[string]string = map[string]string{
	".js":        "\uf2ee",
	".ts":        "\ue69d",
	".json":      "\ueb0f",
	".go":        "\ue627",
	".mod":       "\ue627",
	".sum":       "\ue627",
	".rs":        "\ue7a8",
	".c":         "\ue61e",
	".h":         "\ue61e",
	".cpp":       "\ue61d",
	".cs":        "\ue7b2",
	".html":      "\ue60e",
	".htm":       "\ue60e",
	".css":       "\ue6b8",
	".gitignore": "\ue702",
	".py":        "\ue73c",
	".php":       "\ue73d",
	".java":      "\ue738",
	".rb":        "\ue605",
	".swift":     "\ue699",
	".kt":        "\ue634",
	".yaml":      "\ue8eb",
	".yml":       "\ue8eb",
	".sh":        "\ue691",
	".md":        "\ueb1d",
	".jpg":       "\uf03e",
	".jpeg":      "\uf03e",
	".png":       "\uf03e",
	".webp":      "\uf03e",
	".gif":       "\uf0d78",
	".svg":       "\uf0721",
	".pdf":       "\uf1c1",
	".mp3":       "\ue638",
	".wav":       "\ue638",
}

func GetIcon(fileExtension string) string {
	if icon, ok := icons[fileExtension]; ok {
		return icon
	}

	return "\uf15b"
}
