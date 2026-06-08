package domain

type TaskStatus struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
}

type SearchTaskStatusesFilter struct {
	Name *string
}
