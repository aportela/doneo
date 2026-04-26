package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
)

type TaskResponse struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type TasksResponse struct {
	Success bool           `json:"success"`
	Tasks   []TaskResponse `json:"tasks"`
}

func toTaskResponse(task domain.Task) TaskResponse {
	return TaskResponse{
		ID:      task.ID,
		Summary: task.Summary,
	}
}

func SearchTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := TasksResponse{
		Success: true,
		Tasks: []TaskResponse{
			toTaskResponse(domain.Task{ID: "1", Summary: "Task 1"}),
			toTaskResponse(domain.Task{ID: "2", Summary: "Task 2"}),
		},
	}

	json.NewEncoder(w).Encode(resp)

}
