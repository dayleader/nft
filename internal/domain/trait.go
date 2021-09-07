package domain

import (
	"bytes"
	"errors"
	"image/png"
	"strings"
)

const (
	RarenessKindCommon = RarenessKind(0)
	RarenessKindSilver = RarenessKind(1)
	RarenessKindGold   = RarenessKind(2)
)

// TraitID - trait id.
type TraitID string

// Rareness - rareness.
type RarenessKind int

// TraitWrite struct.
type TraitWrite struct {
	Name         string       `json:"name"`
	Group        *GroupRead   `json:"group"`
	Image        []byte       `json:"image"`
	RarenessKind RarenessKind `json:"rareness"`
}

// TraitRead struct.
type TraitRead struct {
	ID TraitID `json:"id"`
	TraitWrite
}

// ToImageLayer - returns image layer.
func (r *TraitRead) ToImageLayer() (*ImageLayer, error) {
	img, err := png.Decode(bytes.NewReader(r.Image))
	if err != nil {
		return nil, err
	}
	return &ImageLayer{
		Image:    img,
		Priotiry: r.Group.Priotiry,
		XPos:     r.Group.XPos,
		YPos:     r.Group.YPos,
	}, nil
}

// ToERC721 - returns ERC721 trait format.
func (r *TraitRead) ToERC721() (*ERC721Trait, error) {
	if len(r.Group.Name) == 0 {
		return nil, errors.New("trait type required")
	}
	if len(r.Name) == 0 {
		return nil, errors.New("trait value required")
	}
	return &ERC721Trait{
		TraitType: strings.ToUpper(r.Group.Name),
		Value:     strings.ToUpper(r.Name),
	}, nil
}

// TraitRepository - provides access to the storage.
type TraitRepository interface {
	Create(trait *TraitWrite) (TraitID, error)
	GetByID(traitID TraitID) (*TraitRead, error)
	IsExistByName(name string) (bool, error)
	GetAll() ([]*TraitRead, error)
	GetByGroupID(groupID GroupID) ([]*TraitRead, error)
}

// TraitService - provides access to the business logic.
type TraitService interface {
	Import(root string) (int, error)
	GetRandomTraits() ([]*TraitRead, error)
}
