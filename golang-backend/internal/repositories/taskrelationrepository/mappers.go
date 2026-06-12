package taskrelationrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(taskRelation domain.TaskRelation) taskRelationDTO {
	return taskRelationDTO{
		TaskID:        taskRelation.TaskID,
		RelatedTaskID: taskRelation.RelatedTaskID,
		RelationType:  uint8(taskRelation.RelationType),
	}
}

func toDomain(taskRelation taskRelationDTO) domain.TaskRelation {
	return domain.TaskRelation{
		TaskID:        taskRelation.TaskID,
		RelatedTaskID: taskRelation.RelatedTaskID,
		RelationType:  domain.TaskRelationType(taskRelation.RelationType),
	}
}
