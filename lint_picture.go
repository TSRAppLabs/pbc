package pbc

import (
	"fmt"
	"image/png"
	"os"
	. "stash.tsrapplabs.com/ut/jsonlint"
	"strings"
)

var imageRoots = []string{
	"background",
	"footer",
	"icon",
	"logo",
	"strip",
	"thumbnail",
}

func imagePathsFor(filepath string) []string {
	result := [3]string{}

	result[0] = filepath + ".png"
	result[1] = filepath + "@2x.png"
	result[2] = filepath + "@3x.png"

	return result[:]
}

func LintImages(path string) (Warning, error) {

	warn := []string{}

	for _, imageRoot := range imageRoots {

		imagePaths := imagePathsFor(imageRoot)

		existingPaths := getExisting(imagePaths)

		/*
			We want a warning when two are present
			when none are present the image is omitted
			when one is present that image is used everywhere
			when three are present all bases are covered.

		*/
		if len(existingPaths) == 2 { //one is missing
			warn = append(warn, fmt.Sprintf("Missing image file(s): %v", strings.Join(findMissing(imagePaths, existingPaths), ",")))
		}

		warn = append(warn, warnDimension(imagePaths...)...)

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

func warnDimension(files ...string) Warning {
	warn := []string{}
	if len(files) == 0 {
		return warn
	}

	width, height, err := getDimensions(files[0])

	if err != nil {
		return warn
	}

	for base, file := range files[1:] {
		mult := base + 2

		w, h, err := getDimensions(file)

		if err == nil {
			if w != width*mult && h != height*mult {
				warn = append(warn, fmt.Sprintf("%v should be (%v,%v) but is (%v, %v)", file, width*mult, height*mult, w, h))
			}
		}
	}

	return warn
}

func getExisting(paths []string) []string {
	result := []string{}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			result = append(result, path)
		}
	}

	return result
}

func findMissing(given []string, actual []string) []string {
	result := []string{}
	for _, path := range given {

		found := false

		for _, e := range actual {
			if e == path {
				found = true
				break
			}
		}

		if !found {
			result = append(result, path)
		}

	}

	return result
}

/*
  Whimsical sounding isn't it?
*/
func allForOneOneForNone(paths []string) bool {
	existant := 0
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			existant += 1
		}
	}

	return existant == 0 || existant == 1 || existant == len(paths)
}
