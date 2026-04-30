package projectpriorityrepository

type projectPriorityDTO struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Index int    `db:"item_index"`
}
