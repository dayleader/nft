package domain

// GeneratorParams struct.
type GeneratorParams struct {
	Width           int
	Length          int
	InputDirectory  string
	OutputDirectory string
	Number          int
}

// ImageGenerator interface.
type ImageGenerator interface {
	GenerateImages() error
}
