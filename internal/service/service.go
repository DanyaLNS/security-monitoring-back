package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"
	"security-monitor/internal/repository"
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

	rawPayload, _ := json.Marshal(req.RawPayload)
	normalizedPayload, _ := json.Marshal(req.NormalizedPayload)

	occurredAt, err := time.Parse(
		time.RFC3339,
		req.OccurredAt,
	)

	if err != nil {
		return err
	}

	event := &domain.Event{
		ULID: generateULID(),

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

	fmt.Print(event)
	return nil
}

func (s *EventService) GetAll(
	ctx context.Context,
) ([]domain.Event, error) {
	return nil, nil
}

func (s *EventService) Delete(
	ctx context.Context,
	ulid string,
) error {
	return nil
}

func generateULID() string {
	t := time.Now()

	id, _ := ulid.New(
		ulid.Timestamp(t),
		rand.Reader,
	)

	return id.String()
}
