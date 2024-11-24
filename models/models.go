package models

import "context"

// CreateRiskRequest is the API contract for the POST /risks endpoint request.
type CreateRiskRequest struct {
	State       RiskState `json:"state"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// A Risk is the data structure for the Risk object. This also serves as the API
// contract for the GET /risks and /risk/{id} endpoint respopnses. I recommend
// splitting these out into 2 separate structs to avoid leaky abstractions. Reusing the
// struct here for simplicity.
type Risk struct {
	ID          string    `json:"ID"`
	State       RiskState `json:"state"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

// RiskState is one of: open, closed, accepted, investigating
type RiskState string

const (
	StateOpen          RiskState = "open"
	StateClosed        RiskState = "closed"
	StateAccepted      RiskState = "accepted"
	StateInvestigating RiskState = "investigating"
)

//go:generate mockgen . DAOInterface -destination=./mocks/mocks.go
type DAOInterface interface {
	GetAllRisks(ctx context.Context) []Risk
	SaveRisk(ctx context.Context, state RiskState, title, description string) (*Risk, error)
	GetRiskByID(ctx context.Context, id string) (*Risk, error)
}
