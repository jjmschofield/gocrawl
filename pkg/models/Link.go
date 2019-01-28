package models

import "net/url"

type Link struct {
	Id string
	SrcUrl url.URL
	OutUrl url.URL
}