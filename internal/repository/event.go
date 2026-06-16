package repository

import (
	"context"

	"security-monitor/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, e domain.Event) error
	GetAll(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, ulid string) error
}

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, e domain.Event) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO events (
			ulid,
			source_ulid,
			type_ulid,
			severity,
			status,
			title,
			source_ip,
			destination_ip,
			hostname,
			occurred_at,
			raw_payload,
			normalized_payload
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
		)`,
		e.ULID,
		e.SourceULID,
		e.TypeULID,
		e.Severity,
		e.Status,
		e.Title,
		e.SourceIP,
		e.DestinationIP,
		e.Hostname,
		e.OccurredAt,
		e.RawPayload,
		e.NormalizedPayload,
	)

	return err
}

func (r *EventRepo) GetAll(ctx context.Context) ([]domain.Event, error) {
	rows, err := r.db.Query(ctx,
		`SELECT
			ulid,
			source_ulid,
			type_ulid,
			severity,
			status,
			title,
			source_ip,
			destination_ip,
			hostname,
			occurred_at,
			raw_payload,
			normalized_payload,
			created_at,
			updated_at
		FROM events
		ORDER BY occurred_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]domain.Event, 0)

	for rows.Next() {
		var e domain.Event

		if err := rows.Scan(
			&e.ULID,
			&e.SourceULID,
			&e.TypeULID,
			&e.Severity,
			&e.Status,
			&e.Title,
			&e.SourceIP,
			&e.DestinationIP,
			&e.Hostname,
			&e.OccurredAt,
			&e.RawPayload,
			&e.NormalizedPayload,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepo) Delete(ctx context.Context, ulid string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM events WHERE ulid = $1`, ulid)
	return err
}
