package domain

type TaskRelationType uint8

const (
	RelationTypeNone TaskRelationType = iota
	RelationTypeLink
	RelationTypeChild
)

type TaskRelation struct {
	TaskID        string
	RelatedTaskID string
	RelationType  TaskRelationType
}
