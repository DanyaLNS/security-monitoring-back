package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"
	"security-monitor/internal/repository"

	"github.com/oklog/ulid/v2"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(
	repo repository.EventRepository,
) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Create(
	ctx context.Context,
	req dto.CreateEventRequest,
) error {
	if req.Severity < 0 || req.Severity > 10 {
		return fmt.Errorf("severity must be between 0 and 10")
	}

	rawPayload, err := json.Marshal(req.RawPayload)
	if err != nil {
		return err
	}

	var normalizedPayload []byte
	if req.NormalizedPayload != nil {
		normalizedPayload, err = json.Marshal(req.NormalizedPayload)
		if err != nil {
			return err
		}
	}

	occurredAt, err := time.Parse(time.RFC3339, req.OccurredAt)
	if err != nil {
		return err
	}

	ulid, err := generateULID()
	if err != nil {
		return err
	}

	event := domain.Event{
		ULID:              ulid,
		SourceULID:        req.SourceULID,
		TypeULID:          req.TypeULID,
		Severity:          req.Severity,
		Status:            "new",
		Title:             req.Title,
		SourceIP:          req.SourceIP,
		DestinationIP:     req.DestinationIP,
		Hostname:          req.Hostname,
		OccurredAt:        occurredAt,
		RawPayload:        rawPayload,
		NormalizedPayload: normalizedPayload,
	}

	return s.repo.Create(ctx, event)
}

func (s *EventService) GetAll(
	ctx context.Context,
) ([]domain.Event, error) {
	return s.repo.GetAll(ctx)
}

func (s *EventService) Delete(
	ctx context.Context,
	ulid string,
) error {
	return s.repo.Delete(ctx, ulid)
}

func generateULID() (string, error) {
	id, err := ulid.New(
		ulid.Timestamp(time.Now()),
		rand.Reader,
	)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}
