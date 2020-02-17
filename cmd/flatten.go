package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/micahyoung/docker-layer-extract/extract"

	"github.com/urfave/cli"
)

func (b *Builder) FlattenAction(c *cli.Context) error {
	var err error
	imagePath := c.GlobalString("imagefile")
	startlayerID := c.String("startlayerid")
	endlayerID := c.String("endlayerid")
	layerPath := c.String("layerfile")
	stripPax := c.Bool("strip-pax")

	if imagePath == "" {
		return errors.New("missing input image file")
	}
	if layerPath == "" {
		return errors.New("missing output layer file")
	}

	if startlayerID == "" {
		return errors.New("missing desired start layer IDs")
	}

	if _, err = os.Stat(layerPath); !os.IsNotExist(err) {
		return fmt.Errorf("Refusing to overwrite existing file: %s", layerPath)
	}

	allLayerInfos, err := b.extractor.GetImageLayerInfos(imagePath)
	if err != nil {
		return err
	}

	var flattenImageTarballLayerPaths []string
	foundStartLayer := false
	foundEndLayer := false
	for _, layerInfo := range allLayerInfos {
		if layerInfo.ID == startlayerID {
			foundStartLayer = true
		}

		if endlayerID != "" && layerInfo.ID == endlayerID {
			foundEndLayer = true
		}

		if !foundStartLayer {
			continue
		}

		flattenImageTarballLayerPaths = append(flattenImageTarballLayerPaths, layerInfo.LayerPath)

		if foundEndLayer {
			break
		}
	}

	if !foundStartLayer {
		return fmt.Errorf("Desired start layer not found for: %s", startlayerID)
	}

	if endlayerID != "" && !foundEndLayer {
		return fmt.Errorf("Desired end layer not found for: %s", endlayerID)
	}

	flattenerOptions := &extract.FlattenerOptions{
		StripPax: stripPax,
	}

	err = b.flattener.FlattenLayersToPath(imagePath, flattenImageTarballLayerPaths, layerPath, flattenerOptions)
	if err != nil {
		return err
	}

	return nil
}
