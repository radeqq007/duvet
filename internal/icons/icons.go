package icons

var icons map[string]string = map[string]string {
	".js": "\uf2ee",
	".ts": "\ue69d",
	".json": "\ueb0f",
	".go": "\ue627",
}

func GetIcon(fileExtension string) string {
	if icon, ok := icons[fileExtension]; ok {
		return icon
	}

	return "\uf15b"
}

