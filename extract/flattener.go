package extract

import (
	"archive/tar"
	"io"
	"os"
	"path"
	"strings"

	"github.com/micahyoung/docker-layer-extract/layer"
)

type Flattener struct {
	imageRepo         *ImageRepo
	manifestParser    *ManifestParser
	imageConfigParser *ImageConfigParser
	layerReformatter  *layer.LayerReformatter
}

type FlattenerOptions struct {
	StripPax bool
}

func NewFlattener(imageRepo *ImageRepo, manifestParser *ManifestParser, imageConfigParser *ImageConfigParser, layerReformatter *layer.LayerReformatter) *Flattener {
	return &Flattener{imageRepo, manifestParser, imageConfigParser, layerReformatter}
}

func (f *Flattener) FlattenLayersToPath(imagePath string, imageTarballLayerPaths []string, layerPath string, extractorOptions *FlattenerOptions) error {
	var err error
	var imageFile *os.File
	var flattenLayerFile *os.File

	imageFile, err = os.Open(imagePath)
	if err != nil {
		return err
	}
	defer imageFile.Close()

	flattenLayerFile, err = os.Create(layerPath)
	if err != nil {
		return err
	}
	defer flattenLayerFile.Close()

	tarWriter := tar.NewWriter(flattenLayerFile)

	fileMap := map[string]bool{}

	for i := len(imageTarballLayerPaths) - 1; i >= 0; i = i - 1 {
		imageTarballLayerPath := imageTarballLayerPaths[i]

		imageFile.Seek(0, 0)

		imageTarReader := tar.NewReader(imageFile)

		for {
			var imageFileHeader *tar.Header

			imageFileHeader, err = imageTarReader.Next()

			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}

			if imageFileHeader.Name == imageTarballLayerPath {
				layerTarReader := tar.NewReader(imageTarReader)

				for {
					var layerFileHeader *tar.Header

					layerFileHeader, err = layerTarReader.Next()
					if err == io.EOF {
						break
					}
					if err != nil {
						return err
					}

					layerFilePath := layerFileHeader.Name
					fileKey := strings.ToLower(layerFilePath)
					layerFileBasename := path.Base(layerFilePath)
					if strings.HasPrefix(layerFileBasename, ".wh.") {
						deletedLayerFileKey := strings.ToLower(strings.Replace(layerFilePath, ".wh.", "", 1))

						//Add skip entries for whiteout file and eventual deleted file
						fileMap[deletedLayerFileKey] = true
						fileMap[fileKey] = true

						continue
					}

					if fileMap[fileKey] {
						continue
					}

					newHeader := layerFileHeader

					err = tarWriter.WriteHeader(newHeader)
					if err != nil {
						return err
					}

					_, err = io.Copy(tarWriter, layerTarReader)
					if err != nil {
						return err
					}

					fileMap[fileKey] = true
				}
			}
		}
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	err = flattenLayerFile.Close()
	if err != nil {
		return err
	}

	return nil
}
