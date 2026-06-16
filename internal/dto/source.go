package dto

type CreateEventSource struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
}

type EventSourceResponse struct {
	ULID        string  `json:"ulid"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Status      string  `json:"status"`
	EventsCount int     `json:"events_count"`
	LastSeen    *string `json:"last_seen"`
	CreatedAt   string  `json:"created_at"`
}
