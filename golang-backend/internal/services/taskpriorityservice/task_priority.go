package taskpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
)

type TaskPriorityService interface {
	Add(ctx context.Context, taskPriority domain.TaskPriority) error
	Update(ctx context.Context, taskPriority domain.TaskPriority) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.TaskPriority, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error)
}

type taskPriorityService struct {
	repository taskpriorityrepository.TaskPriorityRepository
}

func NewTaskPriorityService(repository taskpriorityrepository.TaskPriorityRepository) TaskPriorityService {
	return &taskPriorityService{repository: repository}
}

func (s *taskPriorityService) Add(ctx context.Context, taskPriority domain.TaskPriority) error {
	return s.repository.Add(ctx, taskpriorityrepository.DomainToDTO(taskPriority))
}

func (s *taskPriorityService) Update(ctx context.Context, taskPriority domain.TaskPriority) error {
	return s.repository.Update(ctx, taskpriorityrepository.DomainToDTO(taskPriority))
}

func (s *taskPriorityService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *taskPriorityService) Get(ctx context.Context, id string) (domain.TaskPriority, error) {
	taskPriority, err := s.repository.Get(ctx, id)
	if err != nil {
		return taskpriorityrepository.DTOToDomain(taskPriority), fmt.Errorf("[TaskPriorityService] failed to get task priority with ID %s: %w", id, err)
	}
	return taskpriorityrepository.DTOToDomain(taskPriority), nil
}

func (s *taskPriorityService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error) {
	taskPriorities, pagerResult, err := s.repository.Search(ctx, pager, order, taskpriorityrepository.DomainFilterToDTO(filter))
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskPriorityService] failed to search task priorities: %w", err)
	}
	return taskpriorityrepository.DTOArrayToDomainArray(taskPriorities), pagerResult, nil
}
