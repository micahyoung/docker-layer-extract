package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func (b *Builder) ListAction(c *cli.Context) error {
	var err error

	imagePath := c.String("imagefile")

	if imagePath == "" {
		return errors.New("missing image file")
	}

	layerInfos, err := b.extractor.GetImageLayerInfos(imagePath)
	if err != nil {
		return err
	}

	for _, layerInfo := range layerInfos {
		friendlyLayerID := strings.Replace(layerInfo.ID, "sha256:", "", 1)
		fmt.Printf("Layer %d:\n", layerInfo.Index)
		fmt.Printf("  ID: %s\n", friendlyLayerID)
		fmt.Printf("  Command: `%s`\n", layerInfo.Command)
	}

	return nil
}
