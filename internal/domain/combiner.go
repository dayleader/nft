package domain

import (
	"image"
	"image/color"
)

// ImageLayer struct.
type ImageLayer struct {
	Image    image.Image
	Priotiry int
	XPos     int
	YPos     int
}

//BgProperty is background property struct.
type BgProperty struct {
	Width   int
	Length  int
	BgColor color.Color
}

// ImageCombiner interface.
type ImageCombiner interface {
	CombineLayers(layers []*ImageLayer, bgProperty *BgProperty) ([]byte, error)
}
