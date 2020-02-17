package helpers

import (
	"archive/tar"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func DockerBuildCommand(dockerfilePath, imageTag string) *exec.Cmd {
	return exec.Command("docker", "build", "--tag", imageTag, "-f", dockerfilePath, filepath.Dir(dockerfilePath))
}

func DockerInspectCommand(imageTag string) *exec.Cmd {
	return exec.Command("docker", "image", "inspect", imageTag)
}

func DockerImageSaveCommand(imageTag, imagePath string) *exec.Cmd {
	return exec.Command("docker", "image", "save", imageTag, "-o", imagePath)
}

func LayerFileContents(tarFile, filePath, daemonOS string) (result string) {
	if daemonOS == "windows" {
		filePath = "Files" + filePath
	}

	layerFile, err := os.Open(tarFile)
	if err != nil {
		return err.Error()
	}

	tarReader := tar.NewReader(layerFile)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err.Error()
		}

		if header.Name == filePath {
			content, err := ioutil.ReadAll(tarReader)
			if err != nil {
				return err.Error()
			}
			return string(content)
		}
	}

	return "not found"
}

func ParseDaemonOS(inspectJSON string) (string, error) {
	var dockerInspects = []struct {
		Os string
	}{}
	err := json.Unmarshal([]byte(inspectJSON), &dockerInspects)
	if err != nil {
		return "", err
	}

	return dockerInspects[0].Os, nil
}

func ParseInspectThirdLayerID(inspectJSON string) (string, error) {
	var dockerInspects = []struct {
		Id     string
		RootFS struct {
			Layers []string
		}
	}{}
	err := json.Unmarshal([]byte(inspectJSON), &dockerInspects)
	if err != nil {
		return "", err
	}

	layers := dockerInspects[0].RootFS.Layers
	thirdLayerID := layers[2]
	friendlyLayerID := strings.Replace(thirdLayerID, "sha256:", "", 1)

	return friendlyLayerID, nil
}
