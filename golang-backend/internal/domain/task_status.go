package domain

type TaskStatus struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Index    int    `json:"index"`
	HexColor string `json:"hexColor"`
}
