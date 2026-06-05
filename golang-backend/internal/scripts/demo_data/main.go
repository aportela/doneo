package demodatascripts

import (
	"github.com/aportela/doneo/internal/database"
)

func CreateDemoData(database database.Database) {
	userIds := createUsers(database, 32)
	roleIds := createRoles(database)
	projectTypeIds := createProjectTypes(database)
	projectPriorityIds := createProjectPriorities(database)
	projectStatusIds := createProjectStatuses(database)
	taskStatusIds := createTaskStatuses(database)
	taskPriorityIds := createTaskPriorities(database)
	createProjects(database, projectTypeIds, projectPriorityIds, projectStatusIds, userIds, roleIds, taskStatusIds, taskPriorityIds, 32)
}
