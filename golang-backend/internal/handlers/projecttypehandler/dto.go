package projecttypehandler

type addProjectTypeRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type updateProjectTypeRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type projectTypeResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type addProjectTypeResponse struct {
	ProjectType projectTypeResponse `json:"projectType"`
}

type updateProjectTypeResponse struct {
	ProjectType projectTypeResponse `json:"projectType"`
}

type getProjectTypeResponse struct {
	ProjectType projectTypeResponse `json:"projectType"`
}

type searchProjectTypesResponse struct {
	ProjectTypes []projectTypeResponse `json:"projectType"`
}
