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

		traitOptions := strings.Split(traitName, ".")
		traitName = traitOptions[0]
		rarenessKind := domain.RarenessKindCommon
		if len(traitOptions) > 1 {
			switch traitOptions[1] {
			case "silver":
				rarenessKind = domain.RarenessKindSilver
			case "gold":
				rarenessKind = domain.RarenessKindGold
			}
		}

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
			Name:         traitName,
			Group:        foundGroup,
			Image:        imgBytes,
			RarenessKind: rarenessKind,
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
		randomTrait, err := GetRandomTrait(traits)
		if err != nil {
			return nil, err
		}
		randomTraits = append(randomTraits, randomTrait)
	}
	if len(groups) != len(randomTraits) {
		return nil, fmt.Errorf(
			fmt.Sprintf("expected traits size %d but got %d", len(groups), len(randomTraits)),
		)
	}
	return randomTraits, nil
}

func GetRandomTrait(traits []*domain.TraitRead) (*domain.TraitRead, error) {

	pdf, err := getProbabilityDensityVector(traits)
	if err != nil {
		return nil, err
	}

	// get cdf
	len := len(traits)
	cdf := make([]float32, len)
	cdf[0] = pdf[0]
	for i := 1; i < len; i++ {
		cdf[i] = cdf[i-1] + pdf[i]
	}
	random := sample(cdf)
	if !(len > random) {
		return nil, fmt.Errorf(
			fmt.Sprintf("random generated trait index out of range, max size: %d, generated index: %d", len, random),
		)
	}
	return traits[random], nil
}

func getProbabilityDensityVector(traits []*domain.TraitRead) ([]float32, error) {
	var (
		len               = len(traits)
		probabilityVector = make([]float32, len)
	)
	var (
		baseChance   = float32(100/len) / 100
		silverChance = baseChance / 2
		goldChance   = baseChance / 4
	)
	var (
		chanceOffset  float32 = 1.00
		commonCounter         = 0
	)
	for i, t := range traits {
		switch t.RarenessKind {
		case domain.RarenessKindSilver:
			probabilityVector[i] = silverChance
			chanceOffset -= silverChance
		case domain.RarenessKindGold:
			probabilityVector[i] = goldChance
			chanceOffset -= goldChance
		default:
			commonCounter++
		}
	}
	for i, p := range probabilityVector {
		if p == 0 {
			probabilityVector[i] = chanceOffset / float32(commonCounter)
		}
	}
	if err := checkProbabilityVector(probabilityVector); err != nil {
		return nil, err
	}
	return probabilityVector, nil
}

func checkProbabilityVector(vector []float32) error {
	var (
		sum         float32 = 0
		controllSum float32 = 1
	)
	for _, p := range vector {
		sum += p
	}
	if !(sum >= controllSum-0.1 || sum >= controllSum+0.1) {
		return fmt.Errorf("Expected probability vector controll sum %v but got %v", controllSum, sum)
	}
	return nil
}

func sample(cdf []float32) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Float32()
	bucket := 0
	for r > cdf[bucket] {
		bucket++
	}
	return bucket
}
