package dao

import (
	"context"
	"maps"
	"slices"

	"arcticwolf.com/cutler/models"
	"github.com/google/uuid"
)

// LocalCache is an implementation of the DAOInterface interface that stores Risks in memory.
type LocalCache struct {
	store map[string]models.Risk
}

var _ = models.DAOInterface(&LocalCache{})

// GetAllRisks returns all risk objects stored in local memory.
func (s *LocalCache) GetAllRisks(_ context.Context) []models.Risk {
	// Consider storing both a map and a slice if memory usage is not a concern
	// and the performance impact of converting map to slice is a concern.
	return slices.Collect(maps.Values(s.store))
}

// SaveRisk saves a single risk object to memory.
func (s *LocalCache) SaveRisk(_ context.Context, state models.RiskState, title, description string) (*models.Risk, error) {
	r := models.Risk{
		ID:          uuid.New().String(),
		State:       state,
		Title:       title,
		Description: description,
	}
	if s.store == nil {
		s.store = make(map[string]models.Risk)
	}
	s.store[r.ID] = r
	return &r, nil
}

// GetRisksByID returns at most one risk object from memory, identified by it's ID.
// This returns an empty response and no error if the ID cannot be found.
func (s *LocalCache) GetRiskByID(_ context.Context, id string) (*models.Risk, error) {
	r, exists := s.store[id]
	if !exists {
		return nil, nil
	}
	return &r, nil
}
