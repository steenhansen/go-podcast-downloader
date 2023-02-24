package globals

import (
	"strings"
	"sync"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type FaultsCollect struct {
	podErrors map[string]error
}

// Record episode errors so can see what episodes are missing from the server
var Faults = FaultsCollect{podErrors: map[string]error{}}
var flawMutex = sync.RWMutex{}

func (faultsCollect *FaultsCollect) Note(mediaUrl string, err error) {
	flawMutex.Lock()
	faultsCollect.podErrors[mediaUrl] = err
	flawMutex.Unlock()

}

func (faultsCollect *FaultsCollect) Clear() {
	for k := range faultsCollect.podErrors {
		delete(faultsCollect.podErrors, k)
	}
}

func (faultsCollect *FaultsCollect) All() (badFiles string) {
	for _, mediaError := range faultsCollect.podErrors {
		errorLines := mediaError.Error()
		singleLine := strings.ReplaceAll(errorLines, "\n", consts.ERROR_SEPARATOR)
		badFiles = badFiles + "\t" + singleLine + "\n"
	}
	return badFiles + "\n"
}
