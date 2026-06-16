package service

import (
	"context"
	"strings"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"
	"security-monitor/internal/repository"
)

type SourceService struct {
	repo repository.EventSourceRepository
}

func NewSourceService(
	repo repository.EventSourceRepository,
) *SourceService {
	return &SourceService{
		repo: repo,
	}
}

func (s *SourceService) Create(
	ctx context.Context,
	input dto.CreateEventSource,
) (*dto.EventSourceResponse, error) {
	id, err := generateULID()
	if err != nil {
		return nil, err
	}

	source := domain.EventSource{
		ULID:        id,
		Name:        strings.TrimSpace(input.Name),
		Description: input.Description,
	}

	if err := s.repo.Create(ctx, source); err != nil {
		return nil, err
	}

	return s.repo.GetByULID(ctx, source.ULID)
}

func (s *SourceService) GetAll(
	ctx context.Context,
) ([]dto.EventSourceResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *SourceService) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.EventSourceResponse, error) {
	return s.repo.GetByULID(ctx, ulid)
}

func (s *SourceService) Delete(
	ctx context.Context,
	ulid string,
) error {
	return s.repo.Delete(ctx, ulid)
}
