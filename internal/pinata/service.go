package pinata

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"nft/internal/domain"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type service struct {
	params domain.IpfsParams
	client *http.Client
}

const (
	pinFileURL = "https://api.pinata.cloud/pinning/pinFileToIPFS"
)

// NewIpfsService - new pinata service.
func NewIpfsService(params domain.IpfsParams) domain.IpfsService {
	return &service{
		params: params,
		client: http.DefaultClient,
	}
}

func (s *service) UploadImages() error {
	var (
		uploadImageParams = make([]*uploadImageParam, 0)
	)
	counter := 0
	err := filepath.Walk(s.params.InputDirectory, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}
		var (
			fileName = info.Name()
		)
		if filepath.Ext(fileName) != ".png" {
			return nil
		}
		counter++
		uploadImageParams = append(uploadImageParams, &uploadImageParam{
			number:   counter,
			path:     path,
			fileName: fileName,
		})
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Files to upload: %d\n", len(uploadImageParams))
	var (
		wg       sync.WaitGroup
		poolSize = 10
		channel  = make(chan *uploadImageParam, poolSize)
	)
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range channel {
				err := s.uploadImage(p)
				if err != nil {
					panic(err)
				}
			}
		}()
	}
	for _, p := range uploadImageParams {
		channel <- p
	}
	return nil
}

type uploadImageParam struct {
	number   int
	path     string
	fileName string
}

func (s *service) uploadImage(p *uploadImageParam) error {
	var (
		key = strings.TrimSuffix(p.fileName, filepath.Ext(p.fileName))
	)
	// Upload image to ipfs.
	imgBytes, err := ioutil.ReadFile(p.path)
	if err != nil {
		return err
	}
	ipfsImageHash, err := s.pinFile(p.fileName, imgBytes, false)
	if err != nil {
		return err
	}

	// Read image traits.
	traitsBytes, err := ioutil.ReadFile(strings.ReplaceAll(p.path, ".png", ".json"))
	if err != nil {
		return err
	}
	traits := []*domain.ERC721Trait{}
	if err := json.Unmarshal(traitsBytes, &traits); err != nil {
		return err
	}

	// Create image metadata file.
	erc721Metadata := &domain.ERC721Metadata{
		Image:      fmt.Sprintf("ipfs://%s", ipfsImageHash),
		Attributes: traits,
	}
	erc721MetadataBytes, err := json.Marshal(erc721Metadata)
	if err != nil {
		return err
	}
	metadataFile, err := os.Create(fmt.Sprintf("%s/%d.%s.json", s.params.OutputDirectory, p.number, key))
	if err != nil {
		return err
	}
	defer metadataFile.Close()
	_, err = metadataFile.Write(erc721MetadataBytes)
	if err != nil {
		return err
	}
	fmt.Printf("Image successfully uploaded to ipfs: %d, %s \n", p.number, p.fileName)
	return nil
}

func (s *service) pinFile(fileName string, data []byte, wrapWithDirectory bool) (string, error) {
	type pinataResponse struct {
		IPFSHash  string `json:"IpfsHash"`
		PinSize   int    `json:"PinSize"`
		Timestamp string `json:"Timestamp"`
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	if _, err := fileWriter.Write(data); err != nil {
		return "", err
	}

	// wrap with directory.
	if wrapWithDirectory {
		fileWriter, err = bodyWriter.CreateFormField("pinataOptions")
		if err != nil {
			return "", err
		}
		if _, err := fileWriter.Write([]byte(`{"wrapWithDirectory": true}`)); err != nil {
			return "", err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", pinFileURL, bodyBuf)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("pinata_api_key", s.params.APIKey)
	req.Header.Set("pinata_secret_api_key", s.params.SecretKey)

	// Do request.
	var (
		retries = 3
		resp    *http.Response
	)
	for retries > 0 {
		resp, err = s.client.Do(req)
		if err != nil {
			retries -= 1
		} else {
			break
		}
	}
	if resp == nil {
		return "", fmt.Errorf("Failed to upload files to ipfs, err: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errMsg := make([]byte, resp.ContentLength)
		_, _ = resp.Body.Read(errMsg)
		return "", fmt.Errorf("Failed to upload file, response code %d, msg: %s", resp.StatusCode, string(errMsg))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	pinataResp := pinataResponse{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&pinataResp)
	if err != nil {
		return "", fmt.Errorf("Failed to decode json, err: %v", err)
	}
	if len(pinataResp.IPFSHash) == 0 {
		return "", errors.New("Ipfs hash not found in the response body")
	}
	return pinataResp.IPFSHash, nil
}
