package combiner

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"nft/internal/domain"
	"sort"
)

type service struct {
}

func NewBasicImageCombiner() domain.ImageCombiner {
	return &service{}
}

func (s *service) CombineLayers(layers []*domain.ImageLayer, bgProperty *domain.BgProperty) ([]byte, error) {

	// Sort list by position.
	layers = sortByPriotiry(layers)

	// Create image's background.
	bgImg := image.NewRGBA(image.Rect(0, 0, bgProperty.Width, bgProperty.Length))

	// Set the background color.
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{bgProperty.BgColor}, image.Point{}, draw.Src)

	// Looping image layers, higher position -> upper layer.
	for _, img := range layers {

		// Set the image offset.
		offset := image.Pt(img.XPos, img.YPos)

		// Combine the image.
		draw.Draw(bgImg, img.Image.Bounds().Add(offset), img.Image, image.Point{}, draw.Over)
	}

	// Encode image to buffer.
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, bgImg); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func sortByPriotiry(list []*domain.ImageLayer) []*domain.ImageLayer {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Priotiry < list[j].Priotiry
	})
	return list
}
