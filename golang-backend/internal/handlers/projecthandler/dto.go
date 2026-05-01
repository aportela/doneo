package projecthandler

type creatorResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectPriorityResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
}

type projectStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
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
	ID          string                  `json:"id"`
	Key         string                  `json:"key"`
	Summary     string                  `json:"summary"`
	Description string                  `json:"description"`
	CreatedBy   creatorResponse         `json:"createdBy"`
	CreatedAt   int64                   `json:"createdAt"`
	Type        projectTypeResponse     `json:"type"`
	Priority    projectPriorityResponse `json:"priority"`
	Status      projectStatusResponse   `json:"status"`
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
