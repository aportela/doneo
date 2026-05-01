package projectpriorityhandler

type addProjectPriorityRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type updateProjectPriorityRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectPriorityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type addProjectPriorityResponse struct {
	ProjectPriority projectPriorityResponse `json:"projectPriority"`
}

type updateProjectPriorityResponse struct {
	ProjectPriority projectPriorityResponse `json:"projectPriority"`
}

type getProjectPriorityResponse struct {
	ProjectPriority projectPriorityResponse `json:"projectPriority"`
}

type searchProjectPrioritysResponse struct {
	ProjectPriorities []projectPriorityResponse `json:"projectPriorities"`
}
