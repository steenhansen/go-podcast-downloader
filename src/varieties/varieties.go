package varieties

import (
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type VarietiesSet map[string]bool

func (varieties VarietiesSet) AddVariety(filename string) {
	if filename != consts.URL_OF_RSS {
		pieces := strings.Split(filename, ".")
		if len(pieces) > 1 {
			variety := pieces[len(pieces)-1]
			varieties[variety] = true
		}
	}
}

func (varieties VarietiesSet) VarietiesString(separator string) (vString string) {
	for k := range varieties {
		vString = vString + k + " "
	}
	vString = strings.TrimSpace(vString)
	return vString
}
