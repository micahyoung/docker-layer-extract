package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ManifestParser struct{}

type manifest struct {
	Config string
	Layers []string
}

func NewManifestParser() *ManifestParser {
	return &ManifestParser{}
}

func (p *ManifestParser) ImageConfigFilename(manifestBuffer *bytes.Buffer) (string, error) {
	var err error
	var imageID string
	var manifests []manifest

	err = json.Unmarshal(manifestBuffer.Bytes(), &manifests)
	if err != nil {
		return "", err
	}

	if len(manifests) != 1 {
		return "", fmt.Errorf("failed to parse manifest")
	}

	imageID = manifests[0].Config
	if imageID == "" {
		return "", fmt.Errorf("failed to find image Id in manifest")
	}

	return imageID, nil
}

func (p *ManifestParser) LayerPaths(manifestBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var layerPaths []string
	var manifests []manifest

	err = json.Unmarshal(manifestBuffer.Bytes(), &manifests)
	if err != nil {
		return nil, err
	}

	if len(manifests) != 1 {
		return nil, fmt.Errorf("failed to parse manifest")
	}

	manifestLayerPaths := manifests[0].Layers
	if len(manifestLayerPaths) == 0 {
		return nil, fmt.Errorf("failed to find any layers in manifest")
	}

	for _, manifestLayerPath := range manifestLayerPaths {
		tarballLayerPath := strings.Replace(manifestLayerPath, `\`, `/`, 1)
		layerPaths = append(layerPaths, tarballLayerPath)
	}

	return layerPaths, nil
}
