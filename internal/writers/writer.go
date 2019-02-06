package writers

import "github.com/jjmschofield/gocrawl/internal/pages"

//go:generate counterfeiter . Writer
type Writer interface {
	Start(in chan pages.Page)
}