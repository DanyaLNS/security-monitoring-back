package repository

import (
	"context"
	"fmt"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, e domain.Event) error
	GetAll(ctx context.Context, filter dto.EventFilter) ([]domain.Event, error)
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

func (r *EventRepo) GetAll(
	ctx context.Context,
	filter dto.EventFilter,
) ([]domain.Event, error) {
	query := `
		SELECT
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
		WHERE 1 = 1
	`

	args := make([]any, 0)
	argID := 1

	if filter.Severity != nil {
		query += fmt.Sprintf(" AND severity = $%d", argID)
		args = append(args, *filter.Severity)
		argID++
	}

	if filter.TypeULID != nil {
		query += fmt.Sprintf(" AND type_ulid = $%d", argID)
		args = append(args, *filter.TypeULID)
		argID++
	}

	if filter.SourceULID != nil {
		query += fmt.Sprintf(" AND source_ulid = $%d", argID)
		args = append(args, *filter.SourceULID)
		argID++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, *filter.Status)
		argID++
	}

	if filter.From != nil {
		query += fmt.Sprintf(" AND occurred_at >= $%d", argID)
		args = append(args, *filter.From)
		argID++
	}

	if filter.To != nil {
		query += fmt.Sprintf(" AND occurred_at <= $%d", argID)
		args = append(args, *filter.To)
		argID++
	}

	if filter.Hostname != nil {
		query += fmt.Sprintf(" AND hostname = $%d", argID)
		args = append(args, *filter.Hostname)
		argID++
	}

	if filter.SourceIP != nil {
		query += fmt.Sprintf(" AND source_ip = $%d", argID)
		args = append(args, *filter.SourceIP)
		argID++
	}

	query += " ORDER BY occurred_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
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
