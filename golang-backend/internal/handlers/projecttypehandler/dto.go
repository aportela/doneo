package projecttypehandler

type addProjectTypeRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type updateProjectTypeRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type ProjectTypeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type addProjectTypeResponse struct {
	ProjectType ProjectTypeResponse `json:"projectType"`
}

type updateProjectTypeResponse struct {
	ProjectType ProjectTypeResponse `json:"projectType"`
}

type getProjectTypeResponse struct {
	ProjectType ProjectTypeResponse `json:"projectType"`
}

type searchProjectTypesResponse struct {
	ProjectTypes []ProjectTypeResponse `json:"projectTypes"`
}
