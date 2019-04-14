package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/micahyoung/docker-layer-extract/extract"

	"github.com/urfave/cli"
)

func (b *Builder) ExtractAction(c *cli.Context) error {
	var err error
	imagePath := c.GlobalString("imagefile")
	layerID := c.String("layerid")
	useNewestLayer := c.Bool("newest")
	layerPath := c.String("layerfile")
	stripPax := c.Bool("strip-pax")

	if imagePath == "" {
		return errors.New("missing input image file")
	}
	if layerPath == "" {
		return errors.New("missing output layer file")
	}

	if layerID == "" && !useNewestLayer {
		return errors.New("missing desired layer ID")
	}

	if _, err = os.Stat(layerPath); !os.IsNotExist(err) {
		return fmt.Errorf("Refusing to overwrite existing file: %s", layerPath)
	}

	layerInfos, err := b.extractor.GetImageLayerInfos(imagePath)
	if err != nil {
		return err
	}

	var imageTarballLayerPath string
	if useNewestLayer {
		newestLayerInfo := layerInfos[len(layerInfos)-1]
		imageTarballLayerPath = newestLayerInfo.LayerPath
	} else {
		for _, layerInfo := range layerInfos {
			if layerInfo.ID == layerID {
				imageTarballLayerPath = layerInfo.LayerPath
			}
		}
	}

	if imageTarballLayerPath == "" {
		return fmt.Errorf("Layer file not found for: %s", layerID)
	}

	extractOptions := &extract.ExtractorOptions{
		StripPax: stripPax,
	}

	err = b.extractor.ExtractLayerToPath(imagePath, imageTarballLayerPath, layerPath, extractOptions)
	if err != nil {
		return err
	}

	return nil
}
