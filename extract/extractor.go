package extract

import (
	"bytes"
	"os"
)

type Extractor struct {
	parser    *Parser
	imageRepo *ImageRepo
}

type layerInfo struct {
	Index     int
	ID        string
	Command   string
	LayerPath string
}

func NewExtractor(parser *Parser, imageRepo *ImageRepo) *Extractor {
	return &Extractor{parser: parser, imageRepo: imageRepo}
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
	imageConfigFilename, err = e.parser.ManifestImageConfigFilename(&manifestBuffer)
	if err != nil {
		return nil, err
	}

	var imageConfigBuffer bytes.Buffer
	err = e.imageRepo.Copy(imagePath, imageConfigFilename, &imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var layerIDs []string
	layerIDs, err = e.parser.ImageConfigLayerIDs(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageCommands []string
	imageCommands, err = e.parser.ImageHistoryCommands(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageTarballLayerPaths []string
	imageTarballLayerPaths, err = e.parser.ManifestLayerPaths(&manifestBuffer)
	if err != nil {
		return nil, err
	}

	for index, layerID := range layerIDs {
		layerInfos = append(layerInfos, &layerInfo{
			Index:     index,
			ID:        layerID,
			Command:   imageCommands[index],
			LayerPath: imageTarballLayerPaths[index],
		})
	}

	return layerInfos, nil
}

func (e *Extractor) ExtractLayerToPath(imagePath, imageTarballLayerPath, layerPath string) error {
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
