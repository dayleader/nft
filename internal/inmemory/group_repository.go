package inmemory

import (
	"errors"
	"fmt"
	"nft/internal/domain"

	"github.com/google/uuid"
)

type groupRepo struct {
	m map[domain.GroupID]*domain.GroupRead
	n map[string]*domain.GroupRead
}

func NewGroupRepository() domain.GroupRepository {
	return &groupRepo{
		m: make(map[domain.GroupID]*domain.GroupRead),
		n: make(map[string]*domain.GroupRead),
	}
}

func (r *groupRepo) Create(group *domain.GroupWrite) (domain.GroupID, error) {
	if group == nil {
		return "", errors.New("cannot store nil group")
	}
	if len(group.Name) == 0 {
		return "", errors.New("cannot store group with empty name")
	}
	groupID := domain.GroupID(uuid.New().String())
	readGroup := &domain.GroupRead{
		ID:         groupID,
		GroupWrite: *group,
	}
	r.m[groupID] = readGroup
	r.n[group.Name] = readGroup
	return groupID, nil
}

func (r *groupRepo) GetByID(groupID domain.GroupID) (*domain.GroupRead, error) {
	foundGroup, ok := r.m[groupID]
	if !ok {
		return nil, fmt.Errorf("group not found by id %v", groupID)
	}
	return foundGroup, nil
}

func (r *groupRepo) GetByName(name string) (*domain.GroupRead, error) {
	foundGroup, ok := r.n[name]
	if !ok {
		return nil, fmt.Errorf("group not found by name %s", name)
	}
	return foundGroup, nil
}

func (r *groupRepo) GetAll() ([]*domain.GroupRead, error) {
	list := make([]*domain.GroupRead, 0)
	for _, group := range r.m {
		list = append(list, group)
	}
	return list, nil
}
