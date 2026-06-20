package domain

type ProjectPriority struct {
	ID       string
	Name     string
	HexColor string
	Index    uint8
}

type SearchProjectPrioritiesFilter struct {
	Name *string
}
