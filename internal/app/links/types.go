package links

import "net/url"

type LinkType string

const (
	InternalPageType LinkType = "internal"
	ExternalPagType  LinkType = "external"
	TelType          LinkType = "tel"
	MailtoType       LinkType = "mailto"
	UnknownType      LinkType = "unknown"
)

func calcType(fromUrl url.URL, toUrl url.URL) LinkType {
	if toUrl.Scheme == "tel" {
		return TelType
	}

	if toUrl.Scheme == "mailto" {
		return MailtoType
	}

	if toUrl.Host == fromUrl.Host && toUrl.Scheme == fromUrl.Scheme {
		return InternalPageType
	} else {
		return ExternalPagType
	}
}
