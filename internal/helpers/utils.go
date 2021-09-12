package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nft/internal/domain"
	"os"
	"path/filepath"
	"strings"
)

func PrintInfo(root string) {
	m := map[string]int{}
	types := map[string]bool{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, errr error) error {
		if filepath.Ext(info.Name()) != ".json" {
			return nil
		}
		jsonBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		traits := []*domain.ERC721Trait{}
		if err := json.Unmarshal(jsonBytes, &traits); err != nil {
			return err
		}
		if len(traits) == 0 {
			return fmt.Errorf("Traits not found for %s", path)
		}
		for _, trait := range traits {
			key := fmt.Sprintf("%s-%s", trait.TraitType, trait.Value)
			m[key] = m[key] + 1
			types[trait.TraitType] = true
		}
		return nil
	})

	for t := range types {
		fmt.Println(t)
		for k, v := range m {
			vals := strings.Split(k, "-")
			if t == vals[0] {
				fmt.Printf("	%s = %d \n", vals[1], v)
			}
		}
	}
	if err != nil {
		panic(err)
	}
}

func PrintInfoV2(dir1, dir2 string) {
	err := filepath.Walk(dir1, func(path string, info os.FileInfo, _ error) error {
		ext := filepath.Ext(info.Name())
		if ext != ".json" {
			return nil
		}
		key := strings.Split(info.Name(), ".")[1]

		// metadata traits
		jsonBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		metadata := domain.ERC721Metadata{}
		if err := json.Unmarshal(jsonBytes, &metadata); err != nil {
			return err
		}
		if len(metadata.Attributes) == 0 {
			return fmt.Errorf("Traits not found for %s", path)
		}

		// traits
		path2 := fmt.Sprintf("%s/%s.json", dir2, key)
		jsonBytes, err = ioutil.ReadFile(path2)
		if err != nil {
			return err
		}
		traits := []*domain.ERC721Trait{}
		if err := json.Unmarshal(jsonBytes, &traits); err != nil {
			return err
		}
		if len(traits) == 0 {
			return fmt.Errorf("Traits not found for %s", path2)
		}
		if err := compareTraits(key, traits, metadata.Attributes); err != nil {
			return err
		}
		return compareTraits(key, metadata.Attributes, traits)
	})
	if err != nil {
		panic(err)
	}
}

func compareTraits(key string, traits1, traits2 []*domain.ERC721Trait) error {
	if len(traits1) == 0 || len(traits2) == 0 {
		return fmt.Errorf("Traits are empty for key %s", key)
	}
	if len(traits1) != len(traits2) {
		return fmt.Errorf("Traits len does not equal, %d vs %d for key %s", len(traits1), len(traits2), key)
	}
	for _, t1 := range traits1 {
		found := false
		for _, t2 := range traits2 {
			if t1 != nil && t2 != nil && *t1 == *t2 {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Trait %v not found for key %s", t1, key)
		}
	}
	return nil
}
