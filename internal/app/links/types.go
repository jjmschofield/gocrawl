package links

import (
	"net/url"
	"regexp"
	"sort"
	"strings"
)

const (
	InternalPageType string = "internal"
	InternalFileType string = "internal_file"
	ExternalPageType string = "external"
	ExternalFileType string = "external_file"
	TelType          string = "tel"
	MailtoType       string = "mailto"
	UnknownType      string = "unknown"
)

// This is a sorted array used in binary search, when adding values add them in alphabetically
var fileExtensions = []string{"asx", "avi", "avi", "doc", "docx", "exe", "f4v", "flv", "gif", "jar", "jar", "jpeg", "jpg", "m1v", "mov", "mp2", "mp4", "mpeg", "mpg", "pdf", "png", "pps", "raw", "rss", "swf", "wav", "wma", "wmv", "xls", "xml", "xsd", "zip"}
var extensionRegex = regexp.MustCompile("\\.[\\w]+$") // .* at end of line

func calcType(fromUrl url.URL, toUrl url.URL) string {
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
	extension := extensionRegex.FindString(testUrl.Path)

	if len(extension) < 1 {
		return false
	}

	extension = strings.TrimLeft(extension, ".")

	fileExtIndex := sort.SearchStrings(fileExtensions, extension)

	// SearchStrings gives us an insert position - so we must test A) we are not out of range and B) we have not found the ext
	if fileExtIndex == len(fileExtensions) || fileExtensions[fileExtIndex] != extension {
		return false
	}

	return true
}
