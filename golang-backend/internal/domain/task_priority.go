package domain

type TaskPriority struct {
	ID       string
	Name     string
	HexColor string
}

type SearchTaskPrioritiesFilter struct {
	Name *string
}
