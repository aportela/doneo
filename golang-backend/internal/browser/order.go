package browser

type Order struct {
	Field string
	Sort  string
}

func (p Order) hasAscendingSort() bool {
	return p.Sort == "ASC"
}

func (p Order) hasDescendingSort() bool {
	return p.Sort == "DESC"
}
