package demodatascripts

import (
	"github.com/aportela/doneo/internal/database"
)

func CreateDemoData(db database.Database) {
	userIds := createUsers(db, 32)
	roleIds := createRoles(db)
	projectTypeIds := createProjectTypes(db)
	projectPriorityIds := createProjectPriorities(db)
	projectStatusIds := createProjectStatuses(db)
	taskStatusIds := createTaskStatuses(db)
	taskPriorityIds := createTaskPriorities(db)
	createProjects(db, projectTypeIds, projectPriorityIds, projectStatusIds, userIds, roleIds, taskStatusIds, taskPriorityIds, 32)
}
