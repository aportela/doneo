package domain

type ProjectPriority struct {
	ID       string
	Name     string
	HexColor string
}

type SearchProjectPrioritiesFilter struct {
	Name *string
}
