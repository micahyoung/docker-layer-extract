package extract

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/micahyoung/docker-layer-extract/layer"
)

type ImageRepo struct {
	layerAnalyzer *layer.LayerReformatter
}

func NewImageRepo(layerAnalyzer *layer.LayerReformatter) *ImageRepo {
	return &ImageRepo{layerAnalyzer}
}

func (i *ImageRepo) Copy(imagePath, filename string, writer io.Writer) error {
	var err error
	var file *os.File

	file, err = os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		var header *tar.Header

		header, err = tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == filename {
			_, err = io.Copy(writer, tarReader)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("%s not found", filename)
}

