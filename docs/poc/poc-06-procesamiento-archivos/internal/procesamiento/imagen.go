package procesamiento

import (
	"github.com/disintegration/imaging"
)

const anchoMaximo = 1600

func ProcesarImagen(inputPath, outputPath string) error {
	img, err := imaging.Open(inputPath, imaging.AutoOrientation(true))
	if err != nil {
		return err
	}

	if img.Bounds().Dx() > anchoMaximo {
		img = imaging.Resize(img, anchoMaximo, 0, imaging.Lanczos)
	}

	return imaging.Save(img, outputPath, imaging.JPEGQuality(80))
}
