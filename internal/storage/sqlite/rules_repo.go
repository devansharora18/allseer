package sqlite

import "allseer/internal/rules"

type RulesRepository struct {
	store *Store
}

func NewRulesRepository(store *Store) *RulesRepository {
	return &RulesRepository{store: store}
}

func (r *RulesRepository) List() ([]rules.Rule, error) {
	return nil, nil
}
