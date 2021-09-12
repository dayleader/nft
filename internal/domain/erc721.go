package domain

// ERC721Trait - ERC721 trait format.
type ERC721Trait struct {
	TraitType string `json:"trait_type"`
	Value     string `json:"value"`
}

// ERC721Metadata - metadata schema.
type ERC721Metadata struct {
	Image      string         `json:"image"`
	Attributes []*ERC721Trait `json:"attributes"`
}
