package pbc

import (
	. "github.com/tsrapplabs/jsonlint"
	"image/png"
	"os"
)

var imageRoots = []string{
	"background",
	"footer",
	"icon",
	"logo",
	"strip",
	"thumbnail",
}

const (
	singleExt = ".png"
	doubleExt = "@2x.png"
	tripleExt = "@3x.png"
)

type CrossPlatformImage struct {
	Root   string
	Single string
	Double string
	Triple string
}

type DimensionExpectation struct {
	Width  int
	Height int
}

func CrossPlatformImageForRoot(name string) CrossPlatformImage {
	var result CrossPlatformImage

	result.Root = name

	if _, err := os.Stat(name + singleExt); err == nil {
		result.Single = name + singleExt
	}

	if _, err := os.Stat(name + doubleExt); err == nil {
		result.Double = name + doubleExt
	}

	if _, err := os.Stat(name + tripleExt); err == nil {
		result.Triple = name + tripleExt
	}

	return result
}

func LintImages(path string) (Warning, error) {

	warn := []string{}

	for _, imageRoot := range imageRoots {

		img := CrossPlatformImageForRoot(imageRoot)

		if img.Single == "" && img.Double == "" && img.Triple == "" {
			continue
		}

		if img.Single == "" {
			warn = append(warn, NewWarning("Missing base image: %v", img.Single)...)
			continue
		}

		warn = append(warn, img.lintSizeRelations()...)
		warn = append(warn, img.lintDimensions()...)
	}

	return warn, nil
}

func getDimensions(filepath string) (int, int, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return 0, 0, err
	}

	img, err := png.Decode(file)

	if err != nil {
		return 0, 0, err
	}

	rec := img.Bounds()

	return rec.Max.X - rec.Min.X, rec.Max.X - rec.Min.Y, nil
}

/*
  Assumption: Single exists
*/
func (img CrossPlatformImage) lintSizeRelations() Warning {
	warn := []string{}
	w, h, err := getDimensions(img.Single)

	if err != nil {
		return warn
	}

	if w2, h2, err := getDimensions(img.Double); err == nil && !(2*w == w2 && 2*h == h2) {
		warn = append(warn, NewWarning("File: %v, expected (%v, %v) but got (%v, %v)", img.Double, 2*w, 2*h, w2, h2)...)
	}

	if w3, h3, err := getDimensions(img.Triple); err == nil && !(3*w == w3 && 3*h == h3) {
		warn = append(warn, NewWarning("File: %v, expected (%v, %v) but got  (%v, %v)", img.Triple, 3*w, 3*h, w3, h3)...)
	}

	return warn
}

func (img CrossPlatformImage) lintDimensions() Warning {
	dimens, ok := dimensions[img.Root]
	warn := []string{}

	if !ok {
		return warn
	}

	w, h, err := getDimensions(img.Single)

	if err == nil && !(dimens.Width == w && dimens.Height == h) {
		warn = append(warn, NewWarning("%v expected dimension (%v, %v) but got (%v, %v) which will result in scaling and cropping", img.Single, dimens.Width, dimens.Height, w, h)...)
	}

	return warn
}

/*
 This tool does not support linting the strip image
*/
var dimensions = map[string]DimensionExpectation{
	"background": DimensionExpectation{Width: 180, Height: 220},
	"footer":     DimensionExpectation{Width: 286, Height: 15},
	"icon":       DimensionExpectation{Width: 60, Height: 60},
	"logo":       DimensionExpectation{Width: 160, Height: 50},
	"thumbnail":  DimensionExpectation{Width: 90, Height: 90},
}
