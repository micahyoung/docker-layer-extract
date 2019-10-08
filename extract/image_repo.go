package extract

import (
	"archive/tar"
	"fmt"
	"io"
	"os"

	"github.com/micahyoung/docker-layer-extract/layer"
)

type ImageRepo struct {
	layerReformatter *layer.LayerReformatter
}

func NewImageRepo() *ImageRepo {
	return &ImageRepo{}
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
