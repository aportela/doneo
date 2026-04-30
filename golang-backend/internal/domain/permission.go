package domain

type Permission int

const (
	Create Permission = 1 << iota
	Update
	Delete
	View
	// TODO: List, Execute, Admin/Full ?
)
