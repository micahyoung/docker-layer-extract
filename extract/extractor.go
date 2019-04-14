package extract

import (
	"bytes"
	"os"

	"github.com/micahyoung/docker-layer-extract/layer"
)

type Extractor struct {
	imageRepo         *ImageRepo
	manifestParser    *ManifestParser
	imageConfigParser *ImageConfigParser
	layerAnalyser     *layer.LayerAnalyser
}

type ExtractorOptions struct {
	StripPax bool
}

type layerInfo struct {
	Index         int
	ID            string
	Command       string
	LayerPath     string
	HasPaxHeaders bool
}

func NewExtractor(imageRepo *ImageRepo, manifestParser *ManifestParser, imageConfigParser *ImageConfigParser, layerAnalyser *layer.LayerAnalyser) *Extractor {
	return &Extractor{imageRepo, manifestParser, imageConfigParser, layerAnalyser}
}

func (e *Extractor) GetImageLayerInfos(imagePath string) ([]*layerInfo, error) {
	var err error
	var layerInfos []*layerInfo

	var manifestBuffer bytes.Buffer
	err = e.imageRepo.Copy(imagePath, "manifest.json", &manifestBuffer)
	if err != nil {
		return nil, err
	}

	var imageConfigFilename string
	imageConfigFilename, err = e.manifestParser.ImageConfigFilename(&manifestBuffer)
	if err != nil {
		return nil, err
	}

	var imageConfigBuffer bytes.Buffer
	err = e.imageRepo.Copy(imagePath, imageConfigFilename, &imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var layerIDs []string
	layerIDs, err = e.imageConfigParser.LayerIDs(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageCommands []string
	imageCommands, err = e.imageConfigParser.HistoryCommands(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageTarballLayerPaths []string
	imageTarballLayerPaths, err = e.manifestParser.LayerPaths(&manifestBuffer)
	if err != nil {
		return nil, err
	}

	for index, layerID := range layerIDs {
		var layerTarBuffer bytes.Buffer
		layerPath := imageTarballLayerPaths[index]
		err = e.imageRepo.Copy(imagePath, layerPath, &layerTarBuffer)
		if err != nil {
			return nil, err
		}

		hasPaxHeaders, err := e.layerAnalyser.LayerHasPaxHeaders(&layerTarBuffer)
		if err != nil {
			return nil, err
		}

		layerInfos = append(layerInfos, &layerInfo{
			Index:         index,
			ID:            layerID,
			Command:       imageCommands[index],
			LayerPath:     layerPath,
			HasPaxHeaders: hasPaxHeaders,
		})
	}

	return layerInfos, nil
}

func (e *Extractor) ExtractLayerToPath(imagePath, imageTarballLayerPath, layerPath string, extractorOptions *ExtractorOptions) error {
	var err error
	var layerFile *os.File

	layerFile, err = os.Create(layerPath)
	if err != nil {
		return err
	}

	err = e.imageRepo.Copy(imagePath, imageTarballLayerPath, layerFile)
	if err != nil {
		return err
	}

	return nil
}
