package writers

import "github.com/jjmschofield/GoCrawl/internal/pages"

//go:generate counterfeiter . Writer
type Writer interface {
	Start(in chan pages.Page)
}