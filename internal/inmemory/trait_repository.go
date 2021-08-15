package inmemory

import (
	"errors"
	"fmt"
	"nft/internal/domain"

	"github.com/google/uuid"
)

type traitRepo struct {
	m map[domain.TraitID]*domain.TraitRead
	n map[string]*domain.TraitRead
	g map[domain.GroupID][]*domain.TraitRead
}

func NewTraitRepository() domain.TraitRepository {
	return &traitRepo{
		m: make(map[domain.TraitID]*domain.TraitRead),
		n: make(map[string]*domain.TraitRead),
		g: make(map[domain.GroupID][]*domain.TraitRead),
	}
}

func (r *traitRepo) Create(trait *domain.TraitWrite) (domain.TraitID, error) {
	if trait == nil {
		return "", errors.New("cannot store nil trait")
	}
	if len(trait.Name) == 0 {
		return "", errors.New("cannot store trait with empty name")
	}
	traitID := domain.TraitID(uuid.New().String())
	traitRead := &domain.TraitRead{
		ID:         traitID,
		TraitWrite: *trait,
	}
	r.m[traitID] = traitRead
	r.n[trait.Name] = traitRead

	_, ok := r.g[trait.Group.ID]
	if !ok {
		r.g[trait.Group.ID] = make([]*domain.TraitRead, 0)
	}
	r.g[trait.Group.ID] = append(r.g[trait.Group.ID], traitRead)
	return traitID, nil
}

func (r *traitRepo) GetByID(traitID domain.TraitID) (*domain.TraitRead, error) {
	foundTrait, ok := r.m[traitID]
	if !ok {
		return nil, fmt.Errorf("trait not found by id %v", traitID)
	}
	return foundTrait, nil
}

func (r *traitRepo) IsExistByName(name string) (bool, error) {
	_, ok := r.n[name]
	return ok, nil
}

func (r *traitRepo) GetAll() ([]*domain.TraitRead, error) {
	list := make([]*domain.TraitRead, 0)
	for _, trait := range r.m {
		list = append(list, trait)
	}
	return list, nil
}

func (r *traitRepo) GetByGroupID(groupID domain.GroupID) ([]*domain.TraitRead, error) {
	return r.g[groupID], nil
}
