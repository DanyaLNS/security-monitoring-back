package repository

import (
	"context"
	"errors"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventTypeRepository interface {
	Create(ctx context.Context, eventType domain.EventType) error
	GetAll(ctx context.Context) ([]dto.EventTypeResponse, error)
	GetByULID(ctx context.Context, ulid string) (*dto.EventTypeResponse, error)
	Delete(ctx context.Context, ulid string) error
}

type EventTypeRepo struct {
	db *pgxpool.Pool
}

func NewEventTypeRepo(db *pgxpool.Pool) *EventTypeRepo {
	return &EventTypeRepo{
		db: db,
	}
}

func (r *EventTypeRepo) Create(
	ctx context.Context,
	eventType domain.EventType,
) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO event_types (
			ulid,
			code,
			name,
			description
		) VALUES (
			$1, $2, $3, $4
		)`,
		eventType.ULID,
		eventType.Code,
		eventType.Name,
		eventType.Description,
	)

	return err
}

func (r *EventTypeRepo) GetAll(
	ctx context.Context,
) ([]dto.EventTypeResponse, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			t.ulid,
			t.code,
			t.name,
			t.description,
			COUNT(e.ulid) AS events_count,
			MAX(e.occurred_at)::text AS last_seen,
			t.created_at::text
		FROM event_types t
		LEFT JOIN events e ON e.type_ulid = t.ulid
		GROUP BY
			t.ulid,
			t.code,
			t.name,
			t.description,
			t.created_at
		ORDER BY
			MAX(e.occurred_at) DESC NULLS LAST,
			t.code ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	eventTypes := make([]dto.EventTypeResponse, 0)

	for rows.Next() {
		var eventType dto.EventTypeResponse

		if err := rows.Scan(
			&eventType.ULID,
			&eventType.Code,
			&eventType.Name,
			&eventType.Description,
			&eventType.EventsCount,
			&eventType.LastSeen,
			&eventType.CreatedAt,
		); err != nil {
			return nil, err
		}

		eventTypes = append(eventTypes, eventType)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return eventTypes, nil
}

func (r *EventTypeRepo) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.EventTypeResponse, error) {
	var eventType dto.EventTypeResponse

	err := r.db.QueryRow(
		ctx,
		`SELECT
			t.ulid,
			t.code,
			t.name,
			t.description,
			COUNT(e.ulid) AS events_count,
			MAX(e.occurred_at)::text AS last_seen,
			t.created_at::text
		FROM event_types t
		LEFT JOIN events e ON e.type_ulid = t.ulid
		WHERE t.ulid = $1
		GROUP BY
			t.ulid,
			t.code,
			t.name,
			t.description,
			t.created_at`,
		ulid,
	).Scan(
		&eventType.ULID,
		&eventType.Code,
		&eventType.Name,
		&eventType.Description,
		&eventType.EventsCount,
		&eventType.LastSeen,
		&eventType.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("event type not found")
		}

		return nil, err
	}

	return &eventType, nil
}

func (r *EventTypeRepo) Delete(
	ctx context.Context,
	ulid string,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM event_types WHERE ulid = $1`,
		ulid,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("event type not found")
	}

	return nil
}
