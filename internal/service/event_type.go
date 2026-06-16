package service

import (
	"context"
	"strings"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"
	"security-monitor/internal/repository"
)

type EventTypeService struct {
	repo repository.EventTypeRepository
}

func NewEventTypeService(
	repo repository.EventTypeRepository,
) *EventTypeService {
	return &EventTypeService{
		repo: repo,
	}
}

func (s *EventTypeService) Create(
	ctx context.Context,
	input dto.CreateEventType,
) (*dto.EventTypeResponse, error) {
	id, err := generateULID()
	if err != nil {
		return nil, err
	}
	eventType := domain.EventType{
		ULID:        id,
		Code:        strings.TrimSpace(input.Code),
		Name:        strings.TrimSpace(input.Name),
		Description: input.Description,
	}

	if err := s.repo.Create(ctx, eventType); err != nil {
		return nil, err
	}

	return s.repo.GetByULID(ctx, eventType.ULID)
}

func (s *EventTypeService) GetAll(
	ctx context.Context,
) ([]dto.EventTypeResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventTypeService) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.EventTypeResponse, error) {
	return s.repo.GetByULID(ctx, ulid)
}

func (s *EventTypeService) Delete(
	ctx context.Context,
	ulid string,
) error {
	return s.repo.Delete(ctx, ulid)
}
