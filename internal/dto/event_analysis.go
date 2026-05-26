package dto

type CreateEventAnalysis struct {
	EventULID       string  `json:"event_ulid" binding:"required"`
	ThreatScore     float64 `json:"threat_score" binding:"required"`
	ThreatLevel     string  `json:"threat_level" binding:"required"`
	DetectionMethod *string `json:"detection_method"`
	Summary         *string `json:"summary"`
	Recommendations *string `json:"recommendations"`
}
