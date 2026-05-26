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

type EventAnalysis struct {
	ULID            string
	EventULID       string
	ThreatScore     float64
	ThreatLevel     string
	DetectionMethod *string
	AnalysisSummary *string
	Recommendations *string
	AnalyzedAt      time.Time
}
