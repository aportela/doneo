package demodatascripts

import (
	"github.com/aportela/doneo/internal/database"
)

func CreateDemoData(database database.Database) {

	userIds := createUsers(database, 32)
	projectTypeIds := createProjectTypes(database)
	createProjects(database, projectTypeIds, userIds, 32)
}
