package cmd

import (
	"github.com/micahyoung/docker-layer-extract/extract"
)

type Builder struct {
	extractor *extract.Extractor
}

func NewBuilder(extractor *extract.Extractor) *Builder {
	return &Builder{extractor: extractor}
}
