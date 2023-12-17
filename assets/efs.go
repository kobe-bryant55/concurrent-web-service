package assets

import (
	"embed"
)

//go:embed "configs"
var EmbeddedFiles embed.FS
