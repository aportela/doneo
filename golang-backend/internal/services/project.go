package services

import (
	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/repositories"
)

type ProjectService struct {
	repo *repositories.ProjectRepository
}

func NewProjectService(repo *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}
func (s *ProjectService) Get(id string) (models.Project, error) {
	return s.repo.GetByID(id)
}
func (s *ProjectService) Search(id string) ([]models.Project, error) {
	return s.repo.Search()
}
