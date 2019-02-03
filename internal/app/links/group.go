package links

type LinkGroup struct {
	Internal     []Link `json:"internal"`
	InternalFile []Link `json:"internalFiles"`
	External     []Link `json:"external"`
	ExternalFile []Link `json:"externalFiles"`
	Tel          []Link `json:"tel"`
	Mailto       []Link `json:"mailto"`
	Unknown      []Link `json:"unknown"`
}

func ToLinkGroup(links []Link) (group LinkGroup){
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
