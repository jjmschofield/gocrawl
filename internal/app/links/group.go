package links

import "encoding/json"

type LinkGroup struct {
	Internal     []Link
	InternalFile []Link
	External     []Link
	ExternalFile []Link
	Tel          []Link
	Mailto       []Link
	Unknown      []Link
}

func ToLinkGroup(links []Link) (group LinkGroup) {
	for _, link := range links {
		switch {
		case link.Type == InternalPageType:
			group.Internal = append(group.Internal, link)
		case link.Type == InternalFileType:
			group.InternalFile = append(group.InternalFile, link)
		case link.Type == ExternalPageType:
			group.External = append(group.External, link)
		case link.Type == ExternalFileType:
			group.ExternalFile = append(group.ExternalFile, link)
		case link.Type == MailtoType:
			group.Mailto = append(group.Mailto, link)
		case link.Type == TelType:
			group.Tel = append(group.Tel, link)
		default:
			group.Unknown = append(group.Unknown, link)
		}
	}

	return group
}

func (group LinkGroup) MarshalJSON() ([]byte, error) {
	basicGroup := struct {
		Internal     []string `json:"internal"`
		InternalFile []string `json:"internalFiles"`
		External     []string `json:"external"`
		ExternalFile []string `json:"externalFiles"`
		Tel          []string `json:"tel"`
		Mailto       []string `json:"mailto"`
		Unknown      []string `json:"unknown"`
	}{
		Internal:     linksToIds(group.Internal),
		InternalFile: linksToIds(group.InternalFile),
		External:     linksToIds(group.External),
		ExternalFile: linksToIds(group.ExternalFile),
		Tel:          linksToIds(group.Tel),
		Mailto:       linksToIds(group.Mailto),
		Unknown:      linksToIds(group.Unknown),
	}

	return json.Marshal(basicGroup)
}

func linksToIds(links []Link) (ids []string) {
	ids = make([]string, len(links))

	for i, link := range links {
		ids[i] = link.Id
	}
	return ids
}
