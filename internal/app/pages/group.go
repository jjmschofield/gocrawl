package pages

import "encoding/json"

type PageGroup struct {
	Internal []Page
}

func ToPageGroup (internal []Page) (group PageGroup){
	group.Internal = append(group.Internal, internal...)
	return group
}

func (group PageGroup) MarshalJSON() ([]byte, error) {
	basicGroup := struct {
		Internal []string `json:"internal"`
	}{
		Internal: pagesToIds(group.Internal),
	}

	return json.Marshal(basicGroup)
}

func pagesToIds(pages []Page) (ids []string) {
	ids = make([]string, len(pages))

	for i, page := range pages {
		ids[i] = page.Id
	}
	return ids
}