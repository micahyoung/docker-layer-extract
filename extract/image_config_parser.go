package extract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type ImageConfigParser struct{}

type imageConfig struct {
	Rootfs struct {
		DiffIDs []string `json:"diff_ids"`
	}
	History []struct {
		CreatedBy string `json:"created_by"`
	}
}

func NewImageConfigParser() *ImageConfigParser {
	return &ImageConfigParser{}
}

func (p *ImageConfigParser) LayerIDs(imageConfigBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var layerIDs []string
	var imageConfig imageConfig

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

func (p *ImageConfigParser) HistoryCommands(imageConfigBuffer *bytes.Buffer) ([]string, error) {
	var err error
	var commands []string
	var imageConfig imageConfig

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
