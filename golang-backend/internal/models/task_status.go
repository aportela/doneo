package models

type TaskStatus struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Index int    `json:"index"`
	Color string `json:"color"`
}
