package taskstatusrepository

type taskStatusDTO struct {
	ID           string `db:"id"`
	Name         string `db:"name"`
	HexColor     string `db:"item_hex_color"`
	Index        uint   `db:"item_index"`
	FlagsBitmask uint64 `db:"flags_bitmask"`
}

type searchFilterDTO struct {
	Name *string
}
