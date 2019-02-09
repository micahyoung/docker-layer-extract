package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/micahyoung/docker-layer-extract/cmd"
	"github.com/micahyoung/docker-layer-extract/extract"
)

func main() {
	app := cli.NewApp()

	extractor := extract.NewExtractor()
	cmdBuilder := cmd.NewBuilder(extractor)

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list layers in image",
			Action:  cmdBuilder.ListAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "imagefile, i",
					Usage: "Image tar file",
				},
			},
		},
		{
			Name:    "extract",
			Aliases: []string{"x"},
			Usage:   "extract image layer",
			Action:  cmdBuilder.ExtractAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "layerid, n",
					Usage: "Layer ID to extract",
				},
				cli.StringFlag{
					Name:  "imagefile, i",
					Usage: "Image tar file",
				},
				cli.StringFlag{
					Name:  "layerfile, o",
					Usage: "Layer tar file",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
