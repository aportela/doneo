package projecthandler

import (
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
)

type creatorResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type addProjectRequest struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Summary string `json:"summary"`
}

type updateProjectRequest struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	Summary string `json:"summary"`
}

type projectResponse struct {
	ID          string                                         `json:"id"`
	Key         string                                         `json:"key"`
	Summary     string                                         `json:"summary"`
	Description string                                         `json:"description"`
	CreatedBy   creatorResponse                                `json:"createdBy"`
	CreatedAt   int64                                          `json:"createdAt"`
	Type        projecttypehandler.ProjectTypeResponse         `json:"type"`
	Priority    projectpriorityhandler.ProjectPriorityResponse `json:"priority"`
	Status      projectstatushandler.ProjectStatusResponse     `json:"status"`
}

type addProjectResponse struct {
	Project projectResponse `json:"project"`
}

type updateProjectResponse struct {
	Project projectResponse `json:"project"`
}

type getProjectResponse struct {
	Project projectResponse `json:"project"`
}

type searchProjectsResponse struct {
	Projects []projectResponse `json:"projects"`
}
