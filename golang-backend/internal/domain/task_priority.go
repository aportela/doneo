package domain

type TaskPriority struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
}

type SearchTaskPrioritiesFilter struct {
	Name *string
}
