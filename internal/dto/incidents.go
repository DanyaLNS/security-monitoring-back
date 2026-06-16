package dto

type CreateIncidentRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description"`
	Severity    int      `json:"severity" binding:"required"`
	Status      *string  `json:"status"`
	EventULIDs  []string `json:"event_ulids"`
}

type UpdateIncidentRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Severity    *int    `json:"severity"`
	Status      *string `json:"status"`
}

type AddIncidentEventsRequest struct {
	EventULIDs []string `json:"event_ulids" binding:"required"`
}

type IncidentResponse struct {
	ULID        string  `json:"ulid"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Severity    int     `json:"severity"`
	Status      string  `json:"status"`
	EventsCount int     `json:"events_count"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
