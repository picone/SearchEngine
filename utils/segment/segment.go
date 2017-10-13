package segment

import (
	"github.com/huichen/sego"
)

var segmenter sego.Segmenter

func init() {
	segmenter.LoadDictionary("/Users/ChienHo/go/src/ChienHo/SearchEngine/data/dictionary.txt")
}

func GetSegmenter() *sego.Segmenter {
	return &segmenter
}
