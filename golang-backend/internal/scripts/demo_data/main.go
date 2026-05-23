package demodatascripts

import (
	"github.com/aportela/doneo/internal/database"
)

func CreateDemoData(database database.Database) {
	userIds := createUsers(database, 32)
	createRoles(database)
	projectTypeIds := createProjectTypes(database)
	projectPriorityIds := createProjectPriorities(database)
	projectStatusIds := createProjectStatuses(database)
	createTaskStatuses(database)
	createTaskPriorities(database)
	createProjects(database, projectTypeIds, projectPriorityIds, projectStatusIds, userIds, 32)
}
