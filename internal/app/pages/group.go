package pages

type PageGroup struct {
	Internal []Page
}

func ToPageGroup (internal []Page) (group PageGroup){
	group.Internal = append(group.Internal, internal...)
	return group
}