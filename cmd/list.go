package cmd

import (
	"errors"
	"fmt"

	"github.com/urfave/cli"
)

func (b *Builder) ListAction(c *cli.Context) error {
	var err error

	imagePath := c.GlobalString("imagefile")

	if imagePath == "" {
		return errors.New("missing image file")
	}

	layerInfos, err := b.extractor.GetImageLayerInfos(imagePath)
	if err != nil {
		return err
	}

	for _, layerInfo := range layerInfos {
		fmt.Printf("Layer %d:\n", layerInfo.Index)
		fmt.Printf("  Command: `%s`\n", layerInfo.Command)
		fmt.Printf("  ID: %s\n", layerInfo.ID)
		fmt.Printf("  ImageLayerPath: %s\n", layerInfo.LayerPath)
	}

	return nil
}
