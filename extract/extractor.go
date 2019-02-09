package extract

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Extractor struct{}

type layerInfo struct {
	Index     int
	ID        string
	Command   string
	LayerPath string
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) GetImageLayerInfos(imagePath string) ([]*layerInfo, error) {
	var err error
	var layerInfos []*layerInfo

	var manifestBuffer bytes.Buffer
	manifestBufferWriter := bufio.NewWriter(&manifestBuffer)
	err = copyTarFileContent(imagePath, "manifest.json", manifestBufferWriter)
	if err != nil {
		return nil, err
	}
	manifestBufferWriter.Flush()

	var imageConfigFilename string
	imageConfigFilename, err = parseManifestImageConfigFilename(&manifestBuffer)
	if err != nil {
		return nil, err
	}

	var imageConfigBuffer bytes.Buffer
	imageConfigBufferWriter := bufio.NewWriter(&imageConfigBuffer)
	err = copyTarFileContent(imagePath, imageConfigFilename, imageConfigBufferWriter)
	if err != nil {
		return nil, err
	}
	imageConfigBufferWriter.Flush()

	var layerIDs []string
	layerIDs, err = parseImageConfigLayerIDs(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageCommands []string
	imageCommands, err = parseImageHistoryCommands(&imageConfigBuffer)
	if err != nil {
		return nil, err
	}

	var imageTarballLayerPaths []string
	imageTarballLayerPaths, err = parseManifestLayerPaths(&manifestBuffer)
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
	err = copyTarFileContent(imagePath, imageTarballLayerPath, layerFile)
	if err != nil {
		return err
	}

	return nil
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

func parseManifestLayerPaths(manifestBuffer *bytes.Buffer) ([]string, error) {
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

func copyTarFileContent(imagePath, filename string, writer io.Writer) error {
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

		// fmt.Printf("%s <> %s\n", header.Name, filename)

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
