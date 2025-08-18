package static

import "embed"

//go:embed css js
var StaticFiles embed.FS

//go:embed templates
var TemplateFiles embed.FS
