package domain

type TimestampFilter struct {
	From   *uint64
	To     *uint64
	Filled *bool
	Empty  *bool
}
