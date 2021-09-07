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

func PrintStatistic(root string) {
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
