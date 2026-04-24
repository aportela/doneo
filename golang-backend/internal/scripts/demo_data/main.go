package demodatascripts

import (
	"github.com/aportela/gotask/internal/database"
)

func CreateDemoData(database database.Database) {

	createUsers(database, 32)
	createProjectTypes(database)
}
