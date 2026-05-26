package repository

import (
	"context"

	"security-monitor/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventSourceRepo struct {
	db *pgxpool.Pool
}

func NewEventSourceRepo(db *pgxpool.Pool) *EventSourceRepo {
	return &EventSourceRepo{db: db}
}

func (r *EventSourceRepo) Create(ctx context.Context, e domain.EventSource) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO event_sources (ulid, name, description)
		 VALUES ($1,$2,$3)`,
		e.ULID, e.Name, e.Description,
	)
	return err
}

func (r *EventSourceRepo) GetAll(ctx context.Context) ([]domain.EventSource, error) {
	rows, err := r.db.Query(ctx, `SELECT ulid, name, description, created_at FROM event_sources`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.EventSource

	for rows.Next() {
		var e domain.EventSource
		rows.Scan(&e.ULID, &e.Name, &e.Description, &e.CreatedAt)
		res = append(res, e)
	}

	return res, nil
}

func (r *EventSourceRepo) Delete(ctx context.Context, ulid string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM event_sources WHERE ulid=$1`, ulid)
	return err
}
