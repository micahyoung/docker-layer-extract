package helpers

import (
	"encoding/json"
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

func ParseInspectLatestLayerID(inspectJSON string) (string, error) {
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
	lastLayerId := layers[len(layers)-1]
	friendlyLayerId := strings.Replace(lastLayerId, "sha256:", "", 1)

	return friendlyLayerId, nil
}
