package repository

import (
	"context"
	"errors"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventSourceRepository interface {
	Create(ctx context.Context, source domain.EventSource) error
	GetAll(ctx context.Context) ([]dto.EventSourceResponse, error)
	GetByULID(ctx context.Context, ulid string) (*dto.EventSourceResponse, error)
	Delete(ctx context.Context, ulid string) error
}

type EventSourceRepo struct {
	db *pgxpool.Pool
}

func NewEventSourceRepo(db *pgxpool.Pool) *EventSourceRepo {
	return &EventSourceRepo{
		db: db,
	}
}

func (r *EventSourceRepo) Create(
	ctx context.Context,
	source domain.EventSource,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO event_sources (
			ulid,
			name,
			description
		) VALUES (
			$1, $2, $3
		)`,
		source.ULID,
		source.Name,
		source.Description,
	)

	return err
}

func (r *EventSourceRepo) GetAll(
	ctx context.Context,
) ([]dto.EventSourceResponse, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			s.ulid,
			s.name,
			s.description,
			CASE
				WHEN COUNT(e.ulid) > 0 THEN 'active'
				ELSE 'inactive'
			END AS status,
			COUNT(e.ulid) AS events_count,
			MAX(e.occurred_at)::text AS last_seen,
			s.created_at::text
		FROM event_sources s
		LEFT JOIN events e ON e.source_ulid = s.ulid
		GROUP BY
			s.ulid,
			s.name,
			s.description,
			s.created_at
		ORDER BY
			MAX(e.occurred_at) DESC NULLS LAST,
			s.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sources := make([]dto.EventSourceResponse, 0)

	for rows.Next() {
		var source dto.EventSourceResponse

		if err := rows.Scan(
			&source.ULID,
			&source.Name,
			&source.Description,
			&source.Status,
			&source.EventsCount,
			&source.LastSeen,
			&source.CreatedAt,
		); err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}

func (r *EventSourceRepo) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.EventSourceResponse, error) {
	var source dto.EventSourceResponse

	err := r.db.QueryRow(
		ctx,
		`SELECT
			s.ulid,
			s.name,
			s.description,
			CASE
				WHEN COUNT(e.ulid) > 0 THEN 'active'
				ELSE 'inactive'
			END AS status,
			COUNT(e.ulid) AS events_count,
			MAX(e.occurred_at)::text AS last_seen,
			s.created_at::text
		FROM event_sources s
		LEFT JOIN events e ON e.source_ulid = s.ulid
		WHERE s.ulid = $1
		GROUP BY
			s.ulid,
			s.name,
			s.description,
			s.created_at`,
		ulid,
	).Scan(
		&source.ULID,
		&source.Name,
		&source.Description,
		&source.Status,
		&source.EventsCount,
		&source.LastSeen,
		&source.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("source not found")
		}

		return nil, err
	}

	return &source, nil
}

func (r *EventSourceRepo) Delete(
	ctx context.Context,
	ulid string,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM event_sources WHERE ulid = $1`,
		ulid,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("source not found")
	}

	return nil
}
