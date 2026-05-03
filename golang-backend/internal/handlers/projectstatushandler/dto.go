package projectstatushandler

type addProjectStatusRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
}

type updateProjectStatusRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
}

type ProjectStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
}

type addProjectStatusResponse struct {
	ProjectStatus ProjectStatusResponse `json:"projectStatus"`
}

type updateProjectStatusResponse struct {
	ProjectStatus ProjectStatusResponse `json:"projectStatus"`
}

type getProjectStatusResponse struct {
	ProjectStatus ProjectStatusResponse `json:"projectStatus"`
}

type searchProjectStatusesRequest struct {
	WorkspaceId string `json:"workspaceId"`
}

type searchProjectStatusesResponse struct {
	ProjectStatuses []ProjectStatusResponse `json:"projectStatuses"`
}
