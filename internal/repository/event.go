package repository

import (
	"context"

	"security-monitor/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetAll(ctx context.Context) ([]domain.Event, error)
	Delete(ctx context.Context, ulid string) error
}

type postgresEventRepository struct {
	db *pgxpool.Pool
}

func NewPostgresEventRepository(db *pgxpool.Pool) EventRepository {
	return &postgresEventRepository{
		db: db,
	}
}

func (r *postgresEventRepository) Create(
	ctx context.Context,
	event *domain.Event,
) error {

	query := `
	INSERT INTO events (
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
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		event.ULID,
		event.SourceULID,
		event.TypeULID,
		event.Severity,
		event.Status,
		event.Title,
		event.SourceIP,
		event.DestinationIP,
		event.Hostname,
		event.OccurredAt,
		event.RawPayload,
		event.NormalizedPayload,
	)

	return err
}

func (r *postgresEventRepository) GetAll(
	ctx context.Context,
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
	ORDER BY occurred_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []domain.Event

	for rows.Next() {
		var e domain.Event

		err := rows.Scan(
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
		)

		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (r *postgresEventRepository) Delete(
	ctx context.Context,
	ulid string,
) error {

	query := `DELETE FROM events WHERE ulid = $1`

	_, err := r.db.Exec(ctx, query, ulid)

	return err
}
