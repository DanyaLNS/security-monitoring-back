package service

import (
	"context"
	"fmt"
	"strings"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"
	"security-monitor/internal/repository"
)

type IncidentService struct {
	repo repository.IncidentRepository
}

func NewIncidentService(
	repo repository.IncidentRepository,
) *IncidentService {
	return &IncidentService{repo: repo}
}

func (s *IncidentService) Create(
	ctx context.Context,
	input dto.CreateIncidentRequest,
) (*dto.IncidentResponse, error) {
	if input.Severity < 0 || input.Severity > 10 {
		return nil, fmt.Errorf("severity must be between 0 and 10")
	}

	status := "open"
	if input.Status != nil && strings.TrimSpace(*input.Status) != "" {
		status = strings.TrimSpace(*input.Status)
	}

	id, err := generateULID()
	if err != nil {
		return nil, err
	}

	incident := domain.Incident{
		ULID:        id,
		Title:       strings.TrimSpace(input.Title),
		Description: input.Description,
		Severity:    input.Severity,
		Status:      status,
	}

	if incident.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if err := s.repo.Create(ctx, incident, input.EventULIDs); err != nil {
		return nil, err
	}

	return s.repo.GetByULID(ctx, incident.ULID)
}

func (s *IncidentService) GetAll(
	ctx context.Context,
) ([]dto.IncidentResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *IncidentService) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.IncidentResponse, error) {
	return s.repo.GetByULID(ctx, ulid)
}

func (s *IncidentService) Update(
	ctx context.Context,
	ulid string,
	input dto.UpdateIncidentRequest,
) (*dto.IncidentResponse, error) {
	if input.Severity != nil {
		if *input.Severity < 0 || *input.Severity > 10 {
			return nil, fmt.Errorf("severity must be between 0 and 10")
		}
	}

	if input.Title != nil {
		trimmed := strings.TrimSpace(*input.Title)
		input.Title = &trimmed

		if trimmed == "" {
			return nil, fmt.Errorf("title cannot be empty")
		}
	}

	if input.Status != nil {
		trimmed := strings.TrimSpace(*input.Status)
		input.Status = &trimmed
	}

	if err := s.repo.Update(ctx, ulid, input); err != nil {
		return nil, err
	}

	return s.repo.GetByULID(ctx, ulid)
}

func (s *IncidentService) Delete(
	ctx context.Context,
	ulid string,
) error {
	return s.repo.Delete(ctx, ulid)
}

func (s *IncidentService) GetEvents(
	ctx context.Context,
	incidentULID string,
) ([]domain.Event, error) {
	return s.repo.GetEvents(ctx, incidentULID)
}

func (s *IncidentService) AddEvents(
	ctx context.Context,
	incidentULID string,
	input dto.AddIncidentEventsRequest,
) (*dto.IncidentResponse, error) {
	if len(input.EventULIDs) == 0 {
		return nil, fmt.Errorf("event_ulids is required")
	}

	if err := s.repo.AddEvents(ctx, incidentULID, input.EventULIDs); err != nil {
		return nil, err
	}

	return s.repo.GetByULID(ctx, incidentULID)
}
