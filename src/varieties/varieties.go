package varieties

import (
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type VarietiesSet map[string]bool

func (varietiesSet VarietiesSet) AddVariety(fileName string) {
	if fileName != consts.URL_OF_RSS_FN {
		varietyPieces := strings.Split(fileName, ".")
		if len(varietyPieces) > 1 {
			variety := varietyPieces[len(varietyPieces)-1]
			varietiesSet[variety] = true
		}
	}
}

func (varietiesSet VarietiesSet) VarietiesString(sepChar string) (allVar string) {
	for aVariety := range varietiesSet {
		allVar = allVar + aVariety + sepChar
	}
	allVar = strings.TrimSpace(allVar)
	return allVar
}
