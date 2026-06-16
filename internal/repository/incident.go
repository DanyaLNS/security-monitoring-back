package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"security-monitor/internal/domain"
	"security-monitor/internal/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository interface {
	Create(ctx context.Context, incident domain.Incident, eventULIDs []string) error
	GetAll(ctx context.Context) ([]dto.IncidentResponse, error)
	GetByULID(ctx context.Context, ulid string) (*dto.IncidentResponse, error)
	Update(ctx context.Context, ulid string, input dto.UpdateIncidentRequest) error
	Delete(ctx context.Context, ulid string) error
	GetEvents(ctx context.Context, incidentULID string) ([]domain.Event, error)
	AddEvents(ctx context.Context, incidentULID string, eventULIDs []string) error
}

type IncidentRepo struct {
	db *pgxpool.Pool
}

func NewIncidentRepo(db *pgxpool.Pool) *IncidentRepo {
	return &IncidentRepo{db: db}
}

func (r *IncidentRepo) Create(
	ctx context.Context,
	incident domain.Incident,
	eventULIDs []string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO incidents (
			ulid,
			title,
			description,
			severity,
			status
		) VALUES (
			$1, $2, $3, $4, $5
		)`,
		incident.ULID,
		incident.Title,
		incident.Description,
		incident.Severity,
		incident.Status,
	)
	if err != nil {
		return err
	}

	for _, eventULID := range eventULIDs {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO incident_events (
				incident_ulid,
				event_ulid
			) VALUES (
				$1, $2
			)
			ON CONFLICT DO NOTHING`,
			incident.ULID,
			eventULID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *IncidentRepo) GetAll(
	ctx context.Context,
) ([]dto.IncidentResponse, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			i.ulid,
			i.title,
			i.description,
			i.severity,
			i.status,
			COUNT(ie.event_ulid) AS events_count,
			i.created_at::text,
			i.updated_at::text
		FROM incidents i
		LEFT JOIN incident_events ie ON ie.incident_ulid = i.ulid
		GROUP BY
			i.ulid,
			i.title,
			i.description,
			i.severity,
			i.status,
			i.created_at,
			i.updated_at
		ORDER BY i.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	incidents := make([]dto.IncidentResponse, 0)

	for rows.Next() {
		var incident dto.IncidentResponse

		if err := rows.Scan(
			&incident.ULID,
			&incident.Title,
			&incident.Description,
			&incident.Severity,
			&incident.Status,
			&incident.EventsCount,
			&incident.CreatedAt,
			&incident.UpdatedAt,
		); err != nil {
			return nil, err
		}

		incidents = append(incidents, incident)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (r *IncidentRepo) GetByULID(
	ctx context.Context,
	ulid string,
) (*dto.IncidentResponse, error) {
	var incident dto.IncidentResponse

	err := r.db.QueryRow(
		ctx,
		`SELECT
			i.ulid,
			i.title,
			i.description,
			i.severity,
			i.status,
			COUNT(ie.event_ulid) AS events_count,
			i.created_at::text,
			i.updated_at::text
		FROM incidents i
		LEFT JOIN incident_events ie ON ie.incident_ulid = i.ulid
		WHERE i.ulid = $1
		GROUP BY
			i.ulid,
			i.title,
			i.description,
			i.severity,
			i.status,
			i.created_at,
			i.updated_at`,
		ulid,
	).Scan(
		&incident.ULID,
		&incident.Title,
		&incident.Description,
		&incident.Severity,
		&incident.Status,
		&incident.EventsCount,
		&incident.CreatedAt,
		&incident.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("incident not found")
		}

		return nil, err
	}

	return &incident, nil
}

func (r *IncidentRepo) Update(
	ctx context.Context,
	ulid string,
	input dto.UpdateIncidentRequest,
) error {
	setParts := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if input.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Severity != nil {
		setParts = append(setParts, fmt.Sprintf("severity = $%d", argID))
		args = append(args, *input.Severity)
		argID++
	}

	if input.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argID))
		args = append(args, *input.Status)
		argID++
	}

	if len(setParts) == 0 {
		return nil
	}

	setParts = append(setParts, "updated_at = NOW()")

	args = append(args, ulid)

	query := fmt.Sprintf(
		`UPDATE incidents
		 SET %s
		 WHERE ulid = $%d`,
		strings.Join(setParts, ", "),
		argID,
	)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("incident not found")
	}

	return nil
}

func (r *IncidentRepo) Delete(
	ctx context.Context,
	ulid string,
) error {
	result, err := r.db.Exec(
		ctx,
		`DELETE FROM incidents WHERE ulid = $1`,
		ulid,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("incident not found")
	}

	return nil
}

func (r *IncidentRepo) GetEvents(
	ctx context.Context,
	incidentULID string,
) ([]domain.Event, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			e.ulid,
			e.source_ulid,
			e.type_ulid,
			e.severity,
			e.status,
			e.title,
			e.source_ip::text,
			e.destination_ip::text,
			e.hostname,
			e.occurred_at,
			e.raw_payload,
			e.normalized_payload,
			e.created_at,
			e.updated_at
		FROM events e
		INNER JOIN incident_events ie ON ie.event_ulid = e.ulid
		WHERE ie.incident_ulid = $1
		ORDER BY e.occurred_at DESC`,
		incidentULID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]domain.Event, 0)

	for rows.Next() {
		var event domain.Event

		if err := rows.Scan(
			&event.ULID,
			&event.SourceULID,
			&event.TypeULID,
			&event.Severity,
			&event.Status,
			&event.Title,
			&event.SourceIP,
			&event.DestinationIP,
			&event.Hostname,
			&event.OccurredAt,
			&event.RawPayload,
			&event.NormalizedPayload,
			&event.CreatedAt,
			&event.UpdatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *IncidentRepo) AddEvents(
	ctx context.Context,
	incidentULID string,
	eventULIDs []string,
) error {
	for _, eventULID := range eventULIDs {
		_, err := r.db.Exec(
			ctx,
			`INSERT INTO incident_events (
				incident_ulid,
				event_ulid
			) VALUES (
				$1, $2
			)
			ON CONFLICT DO NOTHING`,
			incidentULID,
			eventULID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
