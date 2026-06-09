package domain

type ProjectPriority struct {
	ID       string
	Name     string
	HexColor string
	Index    uint
}

type SearchProjectPrioritiesFilter struct {
	Name *string
}
