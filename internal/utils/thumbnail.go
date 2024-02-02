package utils

import (
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"os"
)

func ResizeImage(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	origBounds := img.Bounds()
	origWidth := origBounds.Dx()
	origHeight := origBounds.Dy()

	newWidth := int(float64(origWidth) * 0.8)
	newHeight := int(float64(origHeight) * 0.8)

	dstImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	draw.ApproxBiLinear.Scale(dstImage, dstImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	err = jpeg.Encode(dstFile, dstImage, nil)
	if err != nil {
		return err
	}

	return nil
}
