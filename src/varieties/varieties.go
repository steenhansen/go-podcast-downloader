package varieties

import (
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

// consts?????????
type VarietiesSet map[string]bool

func (varietiesSet VarietiesSet) AddVarietyOLd(fileName string) (variety string) {
	if fileName != consts.URL_OF_RSS_FN {
		varietyPieces := strings.Split(fileName, ".")
		if len(varietyPieces) > 1 {
			variety = varietyPieces[len(varietyPieces)-1]
			varietiesSet[variety] = true
		}
	}
	return variety
}

func (varietiesSet VarietiesSet) VarietiesString(sepChar string) (allVar string) {
	for aVariety := range varietiesSet {
		allVar = allVar + aVariety + sepChar
	}
	allVar = strings.TrimSpace(allVar)
	return allVar
}

func FindVariety(fileName string) (variety string) {
	if fileName != consts.URL_OF_RSS_FN {
		varietyPieces := strings.Split(fileName, ".")
		if len(varietyPieces) > 1 {
			variety = varietyPieces[len(varietyPieces)-1]
		}
	}
	return variety
}

func (varietiesSet VarietiesSet) AddVariety(fileName string) {
	variety := FindVariety(fileName)
	if variety != "" {
		varietiesSet[variety] = true
	}
}
