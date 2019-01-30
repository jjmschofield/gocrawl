package links

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
)

type LinkType string

const (
	InternalPageType LinkType = "internal"
	InternalFileType LinkType = "internal_file"
	ExternalPageType LinkType = "external"
	ExternalFileType LinkType = "external_file"
	TelType          LinkType = "tel"
	MailtoType       LinkType = "mailto"
	UnknownType      LinkType = "unknown"
)

// This is a sorted array used in binary search, when adding values add them in alphabetically
var fileExtensions = []string{"asx", "avi", "avi", "doc", "docx", "exe", "f4v", "flv", "gif", "jar", "jar", "jpeg", "jpg", "m1v", "mov", "mp2", "mp4", "mpeg", "mpg", "pdf", "png", "pps", "raw", "rss", "swf", "wav", "wma", "wmv", "xls", "xml", "xsd", "zip"}

func calcType(fromUrl url.URL, toUrl url.URL) LinkType {
	if toUrl.Scheme == "tel" {
		return TelType
	}

	if toUrl.Scheme == "mailto" {
		return MailtoType
	}

	if toUrl.Host == fromUrl.Host && toUrl.Scheme == fromUrl.Scheme {
		if isFile(toUrl){
			return InternalFileType
		}

		return InternalPageType
	} else {
		if isFile(toUrl){
			return ExternalFileType
		}

		return ExternalPageType
	}
}

// By checking teach path for large number of extensions in a regex we seem to loose quite a bit of time at scale
// Instead we qualify if the url has an extension first and then run test to check if the extension is blacklisted using a binary search
func isFile(testUrl url.URL) bool {
	extensionRegex := regexp.MustCompile("\\.[\\w]+$")
	extension := extensionRegex.FindString(testUrl.Path)

	if len(extension) < 1 {
		return false
	}

	extension = strings.TrimLeft(extension, ".")

	fileExtIndex := sort.SearchStrings(fileExtensions, extension);

	if fileExtIndex < 0 {
		return false
	}

	return true
}
