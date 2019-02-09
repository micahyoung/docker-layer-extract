package extract

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Extractor struct{}

type layerInfo struct {
	Index   int
	ID      string
	Command string
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) GetImageLayerInfos(imagePath string) ([]*layerInfo, error) {
	var err error
	var layerInfos []*layerInfo

	var manifestBuffer *bytes.Buffer
	manifestBuffer, err = getTarFileContent(imagePath, "manifest.json")
	if err != nil {
		return nil, err
	}

	var imageConfigFilename string
	imageConfigFilename, err = parseManifestImageConfigFilename(manifestBuffer)
	if err != nil {
		return nil, err
	}

	var imageConfigBuffer *bytes.Buffer
	imageConfigBuffer, err = getTarFileContent(imagePath, imageConfigFilename)
	if err != nil {
		return nil, err
	}

	var layerIDs []string
	layerIDs, err = parseImageConfigLayerIDs(imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageCommands []string
	imageCommands, err = parseImageHistoryCommands(imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	for index, layerID := range layerIDs {
		layerInfos = append(layerInfos, &layerInfo{
			Index:   index,
			ID:      layerID,
			Command: imageCommands[index],
		})
	}

	return layerInfos, nil
}

func parseManifestImageConfigFilename(manifestBuffer *bytes.Buffer) (string, error) {
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

func parseImageConfigLayerIDs(imageConfigBuffer *bytes.Buffer) ([]string, error) {
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

	layerIDs = imageConfig.Rootfs.DiffIDs
	if len(layerIDs) == 0 {
		return nil, fmt.Errorf("failed to parse manifest")
	}

	return layerIDs, nil
}

func parseImageHistoryCommands(imageConfigBuffer *bytes.Buffer) ([]string, error) {
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

func getTarFileContent(imagePath, filename string) (*bytes.Buffer, error) {
	var err error
	var file *os.File

	file, err = os.Open(imagePath)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if header.Name == filename {
			fileBuffer := new(bytes.Buffer)
			fileBuffer.ReadFrom(tarReader)
			return fileBuffer, nil
		}
	}

	return nil, fmt.Errorf("%s not found", filename)
}
