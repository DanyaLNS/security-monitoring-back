package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourname/security-monitor/internal/domain"
)

type EventAnalysisRepo struct {
	db *pgxpool.Pool
}

func NewEventAnalysisRepo(db *pgxpool.Pool) *EventAnalysisRepo {
	return &EventAnalysisRepo{db: db}
}

func (r *EventAnalysisRepo) Create(ctx context.Context, e domain.EventAnalysis) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO event_analysis (
		 ulid, event_ulid, threat_score, threat_level,
		 detection_method, analysis_summary, recommendations
		) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		e.ULID, e.EventULID, e.ThreatScore, e.ThreatLevel,
		e.DetectionMethod, e.AnalysisSummary, e.Recommendations,
	)
	return err
}
