package domain

import "time"

type Event struct {
	ULID              string
	SourceULID        string
	TypeULID          string
	Severity          int
	Status            string
	Title             string
	SourceIP          *string
	DestinationIP     *string
	Hostname          *string
	OccurredAt        time.Time
	RawPayload        []byte
	NormalizedPayload []byte
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
