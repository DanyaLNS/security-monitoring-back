package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"security-monitor/internal/domain"
)

type EventTypeRepo struct {
	db *pgxpool.Pool
}

func NewEventTypeRepo(db *pgxpool.Pool) *EventTypeRepo {
	return &EventTypeRepo{db: db}
}

func (r *EventTypeRepo) Create(ctx context.Context, e domain.EventType) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO event_types (ulid, code, name, description)
		 VALUES ($1,$2,$3,$4)`,
		e.ULID, e.Code, e.Name, e.Description,
	)
	return err
}

func (r *EventTypeRepo) GetAll(ctx context.Context) ([]domain.EventType, error) {
	rows, err := r.db.Query(ctx, `SELECT ulid, code, name, description, created_at FROM event_types`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []domain.EventType

	for rows.Next() {
		var e domain.EventType
		rows.Scan(&e.ULID, &e.Code, &e.Name, &e.Description, &e.CreatedAt)
		res = append(res, e)
	}

	return res, nil
}

func (r *EventTypeRepo) Delete(ctx context.Context, ulid string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM event_types WHERE ulid=$1`, ulid)
	return err
}
