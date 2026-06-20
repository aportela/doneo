package domain

type TaskPriority struct {
	ID       string
	Name     string
	HexColor string
	Index    uint8
}

type SearchTaskPrioritiesFilter struct {
	Name *string
}
