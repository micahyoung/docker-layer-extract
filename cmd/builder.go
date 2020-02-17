package cmd

import (
	"github.com/micahyoung/docker-layer-extract/extract"
)

type Builder struct {
	extractor *extract.Extractor
	flattener *extract.Flattener
}

func NewBuilder(extractor *extract.Extractor, flattener *extract.Flattener) *Builder {
	return &Builder{extractor, flattener}
}
