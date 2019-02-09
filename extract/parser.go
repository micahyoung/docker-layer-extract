package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ManifestImageConfigFilename(manifestBuffer *bytes.Buffer) (string, error) {
	var err error
	var imageID string
	var manifests = []struct {
		Config string
	}{}

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

func (p *Parser) ManifestLayerPaths(manifestBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var layerPaths []string
	var manifests = []struct {
		Layers []string
	}{}

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

func (p *Parser) ImageConfigLayerIDs(imageConfigBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var layerIDs []string
	var imageConfig = struct {
		Rootfs struct {
			DiffIDs []string `json:"diff_ids"` //fFIXME, remove json if possible
		}
	}{}

	err = json.Unmarshal(imageConfigBuffer.Bytes(), &imageConfig)
	if err != nil {
		return nil, err
	}

	shasumLayerIDs := imageConfig.Rootfs.DiffIDs
	if len(shasumLayerIDs) == 0 {
		return nil, fmt.Errorf("failed to parse manifest")
	}

	for _, shasumLayerID := range shasumLayerIDs {
		friendlyLayerID := strings.Replace(shasumLayerID, "sha256:", "", 1)
		layerIDs = append(layerIDs, friendlyLayerID)
	}

	return layerIDs, nil
}

func (p *Parser) ImageHistoryCommands(imageConfigBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var commands []string
	var imageConfig = struct {
		History []struct {
			CreatedBy string `json:"created_by"`
		}
	}{}

	err = json.Unmarshal(imageConfigBuffer.Bytes(), &imageConfig)
	if err != nil {
		return nil, err
	}

	for _, historyItem := range imageConfig.History {
		commands = append(commands, historyItem.CreatedBy)
	}

	if len(commands) == 0 {
		return nil, fmt.Errorf("failed to parse manifest")
	}

	return commands, nil
}
