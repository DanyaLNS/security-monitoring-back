package dto

type CreateEventRequest struct {
	SourceULID        string      `json:"source_ulid" binding:"required"`
	TypeULID          string      `json:"type_ulid" binding:"required"`
	Severity          int         `json:"severity" binding:"required"`
	Title             string      `json:"title" binding:"required"`
	SourceIP          *string     `json:"source_ip"`
	DestinationIP     *string     `json:"destination_ip"`
	Hostname          *string     `json:"hostname"`
	OccurredAt        string      `json:"occurred_at" binding:"required"`
	RawPayload        interface{} `json:"raw_payload" binding:"required"`
	NormalizedPayload interface{} `json:"normalized_payload"`
}
