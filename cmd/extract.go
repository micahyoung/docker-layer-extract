package cmd

import (
	"errors"

	"github.com/urfave/cli"
)

func (b *Builder) ExtractAction(c *cli.Context) error {
	imagePath := c.String("imagefile")
	layerPath := c.String("layerfile")
	layerID := c.String("layerid")

	if imagePath == "" {
		return errors.New("missing image file")
	}
	if layerPath == "" {
		return errors.New("missing layer file")
	}
	if layerID == "" {
		return errors.New("missing layer ID")
	}

	return nil
}
