package taskrelationrepository

type taskRelationDTO struct {
	TaskID        string `db:"task_id"`
	RelatedTaskID string `db:"related_task_id"`
	RelationType  uint8  `db:"relation_type"`
}
