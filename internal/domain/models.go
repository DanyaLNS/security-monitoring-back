package domain

import "time"

type EventSource struct {
	ULID        string
	Name        string
	Description *string
	CreatedAt   time.Time
}

type EventType struct {
	ULID        string
	Code        string
	Name        string
	Description *string
	CreatedAt   time.Time
}

type Event struct {
	ULID              string    `json:"ulid"`
	SourceULID        string    `json:"source_ulid"`
	TypeULID          string    `json:"type_ulid"`
	Severity          int       `json:"severity"`
	Status            string    `json:"status"`
	Title             string    `json:"title"`
	SourceIP          *string   `json:"source_ip"`
	DestinationIP     *string   `json:"destination_ip"`
	Hostname          *string   `json:"hostname"`
	OccurredAt        time.Time `json:"occurred_at"`
	RawPayload        []byte    `json:"raw_payload"`
	NormalizedPayload []byte    `json:"normalized_payload"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Incident struct {
	ULID        string
	Title       string
	Description *string
	Severity    int
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
