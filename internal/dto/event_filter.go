package dto

type EventFilter struct {
	Severity   *int
	TypeULID   *string
	SourceULID *string
	Status     *string
	From       *string
	To         *string
	Hostname   *string
	SourceIP   *string
}
