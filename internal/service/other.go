package service

import (
	"context"

	"github.com/gin-gonic/gin"
)

type SourceService struct{}
type EventTypeService struct{}
type DashboardService struct{}
type IncidentService struct{}
type AnalysisService struct{}

func (s *SourceService) Create(ctx context.Context, input any) (any, error) {
	return input, nil
}

func (s *SourceService) GetAll(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *SourceService) GetByULID(ctx context.Context, ulid string) (any, error) {
	return map[string]any{"ulid": ulid}, nil
}

func (s *SourceService) Delete(ctx context.Context, ulid string) error {
	return nil
}

func (s *EventTypeService) Create(ctx context.Context, input any) (any, error) {
	return input, nil
}

func (s *EventTypeService) GetAll(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *EventTypeService) GetByULID(ctx context.Context, ulid string) (any, error) {
	return map[string]any{"ulid": ulid}, nil
}

func (s *EventTypeService) Delete(ctx context.Context, ulid string) error {
	return nil
}

func (s *DashboardService) GetMetrics(ctx context.Context) (any, error) {
	return gin.H{}, nil
}

func (s *DashboardService) GetTimeline(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *DashboardService) GetSeverityDistribution(ctx context.Context) (any, error) {
	return gin.H{}, nil
}

func (s *DashboardService) GetTopSources(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *DashboardService) GetRecentEvents(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *IncidentService) GetAll(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *IncidentService) GetByULID(ctx context.Context, ulid string) (any, error) {
	return map[string]any{"ulid": ulid}, nil
}

func (s *IncidentService) GetIncidentEvents(ctx context.Context, ulid string) (any, error) {
	return []any{}, nil
}

func (s *AnalysisService) GetAll(ctx context.Context) (any, error) {
	return []any{}, nil
}

func (s *AnalysisService) GetByEventULID(ctx context.Context, ulid string) (any, error) {
	return map[string]any{"event_ulid": ulid}, nil
}
