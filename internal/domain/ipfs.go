package domain

// IpfsParams - ipfs parameters.
type IpfsParams struct {
	InputDirectory  string
	OutputDirectory string
	APIKey          string
	SecretKey       string
}

// IpfsService - ipfs service.
type IpfsService interface {
	UploadImages() error
}
