package trait

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"nft/internal/domain"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type service struct {
	groupRepository domain.GroupRepository
	traitRepository domain.TraitRepository
}

// NewBasicTraitService - returns a naïve, stateless implementation of a service.
func NewBasicTraitService(
	groupRepository domain.GroupRepository,
	traitRepository domain.TraitRepository,
) domain.TraitService {
	return &service{
		groupRepository: groupRepository,
		traitRepository: traitRepository,
	}
}

func (s *service) Import(root string) (int, error) {
	priority := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, errr error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".png" {
			return fmt.Errorf("Bad file extension: %s", filepath.Ext(info.Name()))
		}
		splited := strings.Split(path, "/")
		groupName := splited[1][3:]
		traitName := info.Name()[:len(info.Name())-4]

		// Create group if not exist.
		foundGroup, _ := s.groupRepository.GetByName(groupName)
		if foundGroup == nil {
			_, err := s.groupRepository.Create(&domain.GroupWrite{
				Name:     groupName,
				Priotiry: priority,
			})
			priority++

			fmt.Printf("%s - %d\n", groupName, priority)
			if err != nil {
				return err
			}
		}
		foundGroup, _ = s.groupRepository.GetByName(groupName)
		if foundGroup == nil {
			return fmt.Errorf("group not found %s", groupName)
		}

		// Read image.
		imgBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		// Create trait.
		_, err = s.traitRepository.Create(&domain.TraitWrite{
			Name:  traitName,
			Group: foundGroup,
			Image: imgBytes,
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (s *service) GetRandomTraits() ([]*domain.TraitRead, error) {
	// Get all groups.
	groups, err := s.groupRepository.GetAll()
	if err != nil {
		return nil, err
	}

	// Group traits by their types.
	groupedTraits := make(map[domain.GroupID][]*domain.TraitRead)
	for _, group := range groups {
		traits, err := s.traitRepository.GetByGroupID(group.ID)
		if err != nil {
			return nil, err
		}
		if len(traits) == 0 {
			return nil, fmt.Errorf("traits not found for group %v (%s)", group.ID, group.Name)
		}
		groupedTraits[group.ID] = traits
	}

	// Choose random traits for each available group.
	randomTraits := make([]*domain.TraitRead, 0)
	for _, traits := range groupedTraits {
		rand.Seed(time.Now().UnixNano())
		var (
			max    = len(traits)
			min    = 0
			random = rand.Intn(max-min) + min
		)
		if !(max > random) {
			return nil, fmt.Errorf(
				fmt.Sprintf("random generated trait index out of range, max size: %d, generated index: %d", max, random),
			)
		}
		randomTraits = append(randomTraits, traits[random])
	}
	if len(groups) != len(randomTraits) {
		return nil, fmt.Errorf(
			fmt.Sprintf("expected traits size %d but got %d", len(groups), len(randomTraits)),
		)
	}
	return randomTraits, nil
}
