package main

import (
	"fmt"
	"nft/internal/domain"
	"nft/internal/trait"
)

func main() {

	// gold = 1
	// silver = 3
	// common = 6
	var traits = []*domain.TraitRead{
		{
			TraitWrite: domain.TraitWrite{
				Name:         "1",
				RarenessKind: domain.RarenessKindGold,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "2",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "3",
				RarenessKind: domain.RarenessKindSilver,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "4",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "5",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "6",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "7",
				RarenessKind: domain.RarenessKindSilver,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "8",
				RarenessKind: domain.RarenessKindSilver,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "9",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
		{
			TraitWrite: domain.TraitWrite{
				Name:         "10",
				RarenessKind: domain.RarenessKindCommon,
			},
		},
	}

	m := map[int]int{}
	for i := 0; i < 100; i++ {
		t, err := trait.GetRandomTrait(traits)
		if err != nil {
			panic(err)
		}
		m[int(t.RarenessKind)] = m[int(t.RarenessKind)] + 1
	}

	fmt.Println(m)
}
