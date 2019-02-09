package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func (b *Builder) ExtractAction(c *cli.Context) error {
	var err error
	imagePath := c.GlobalString("imagefile")
	layerID := c.String("layerid")
	layerPath := c.String("layerfile")

	if imagePath == "" {
		return errors.New("missing input image file")
	}
	if layerID == "" {
		return errors.New("missing desired layer ID")
	}
	if layerPath == "" {
		return errors.New("missing output layer file")
	}

	if _, err = os.Stat(layerPath); !os.IsNotExist(err) {
		return fmt.Errorf("Refusing to overwrite existing file: %s", layerPath)
	}

	layerInfos, err := b.extractor.GetImageLayerInfos(imagePath)
	if err != nil {
		return err
	}

	var imageTarballLayerPath string
	for _, layerInfo := range layerInfos {
		if layerInfo.ID == layerID {
			imageTarballLayerPath = layerInfo.LayerPath
		}
	}

	if imageTarballLayerPath == "" {
		return fmt.Errorf("Layer file not found for: %s", layerID)
	}

	err = b.extractor.ExtractLayerToPath(imagePath, imageTarballLayerPath, layerPath)
	if err != nil {
		return err
	}

	return nil
}
