package projecthandler

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
	ID      string `json:"id"`
	Key     string `json:"key"`
	Summary string `json:"summary"`
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
