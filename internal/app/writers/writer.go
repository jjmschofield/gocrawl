package writers

import "github.com/jjmschofield/GoCrawl/internal/app/pages"

type Writer interface {
	Start(in chan pages.Page)
}