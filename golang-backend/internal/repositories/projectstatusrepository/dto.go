package projectstatusrepository

type projectStatusDTO struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Index uint   `db:"item_index"`
}
