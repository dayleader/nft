package generator

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"image/color"
	"nft/internal/domain"
	"os"
	"sort"
	"sync"
)

type service struct {
	params        domain.GeneratorParams
	traitService  domain.TraitService
	imageCombiner domain.ImageCombiner
}

// NewBasicRarityService returns a naïve, stateless implementation of a service.
func NewBasicImageGenerator(
	params domain.GeneratorParams,
	traitService domain.TraitService,
	imageCombiner domain.ImageCombiner,
) domain.ImageGenerator {
	return &service{
		params:        params,
		traitService:  traitService,
		imageCombiner: imageCombiner,
	}
}

func (s *service) GenerateImages() error {
	if !(s.params.Number > 0) {
		return fmt.Errorf("specify at least one element to generate")
	}
	var (
		wg       sync.WaitGroup
		poolSize = 20
		channel  = make(chan int, poolSize)
	)
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _ = range channel {
				err := s.generateImageInternal()
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	for i := 0; i < s.params.Number; i++ {
		channel <- i
	}
	return nil
}

func (s *service) generateImageInternal() error {
	img, traits, err := s.generateInternal()
	if err != nil {
		return err
	}
	key := ""
	for _, flat := range traits {
		key = key + " " + fmt.Sprintf("%s-%s", flat.TraitType, flat.Value)
	}
	key2 := hash(key)
	// Create image file.
	imgFile, err := os.Create(fmt.Sprintf("%s/%d.png", s.params.OutputDirectory, key2))
	if err != nil {
		return err
	}
	defer imgFile.Close()
	_, err = imgFile.Write(img)
	if err != nil {
		return err
	}
	// Create traits file.
	traitsFile, err := os.Create(fmt.Sprintf("%s/%d.json", s.params.OutputDirectory, key2))
	if err != nil {
		return err
	}
	defer traitsFile.Close()
	traitsBytes, err := json.Marshal(traits)
	if err != nil {
		return err
	}
	_, err = traitsFile.Write(traitsBytes)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) generateInternal() ([]byte, []*domain.ERC721Trait, error) {
	// Get random traits.
	traits, err := s.traitService.GetRandomTraits()
	if err != nil {
		return nil, nil, err
	}
	if len(traits) == 0 {
		return nil, nil, errors.New("cannot generate with empty random traits")
	}

	// Convert traints to ERC721 format.
	erc721Traits := make([]*domain.ERC721Trait, len(traits))
	for i, trait := range traits {
		erc721Trait, err := trait.ToERC721()
		if err != nil {
			return nil, nil, err
		}
		erc721Traits[i] = erc721Trait
	}

	// Convert traits to image layers.
	layers := make([]*domain.ImageLayer, len(traits))
	for i, trait := range traits {
		layer, err := trait.ToImageLayer()
		if err != nil {
			return nil, nil, err
		}
		layers[i] = layer
	}

	// Combine image layers together.
	img, err := s.imageCombiner.CombineLayers(layers, &domain.BgProperty{
		Width:   s.params.Width,
		Length:  s.params.Length,
		BgColor: color.Transparent,
	})
	if err != nil {
		return nil, nil, err
	}
	return img, sortTraits(erc721Traits), nil
}

func sortTraits(list []*domain.ERC721Trait) []*domain.ERC721Trait {
	sort.Slice(list, func(i, j int) bool {
		t1 := hash(list[i].TraitType)
		t2 := hash(list[j].TraitType)
		return t1 < t2
	})
	return list
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
