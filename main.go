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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "imagefile, i",
			Usage: "Image tar file (get from: docker save)",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list layers in image",
			Action:  cmdBuilder.ListAction,
		},
		{
			Name:    "extract",
			Aliases: []string{"x"},
			Usage:   "extract image layer",
			Action:  cmdBuilder.ExtractAction,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "layerfile, o",
					Usage: "Output layer tar file",
				},
				cli.StringFlag{
					Name:  "layerid, l",
					Usage: "Layer ID to extract (get from: docker-layer-extract list)",
				},
				cli.BoolFlag{
					Name:  "newest, n",
					Usage: "Use the most recent layer",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
