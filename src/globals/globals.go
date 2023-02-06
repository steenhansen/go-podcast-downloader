package globals

import (
	"fmt"
	"sync"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type FaultsCollect struct {
	podErrors map[string]error
}

var Faults = FaultsCollect{podErrors: map[string]error{}}

func (faultsCollect *FaultsCollect) Note(mediaUrl string, err error) {
	var mu sync.Mutex
	mu.Lock()
	faultsCollect.podErrors[mediaUrl] = err
	mu.Unlock()

}

func (faultsCollect *FaultsCollect) Clear() {
	for k := range faultsCollect.podErrors {
		delete(faultsCollect.podErrors, k)
	}
}

func (faultsCollect *FaultsCollect) All() (badFiles string) {
	for _, mediaError := range faultsCollect.podErrors {
		badFiles = badFiles + "\t\t" + mediaError.Error() + "\n"
	}
	return badFiles
}

type ConsoleCollect struct {
	progText string
}

var Console = ConsoleCollect{progText: ""}

func (consoleCollect *ConsoleCollect) Note(progressStr string) {
	fmt.Print(progressStr)
	if progressStr != consts.CLEAR_SCREEN {
		consoleCollect.progText = consoleCollect.progText + progressStr
	}
}

func (consoleCollect *ConsoleCollect) Clear() {
	consoleCollect.progText = ""
}

func (consoleCollect *ConsoleCollect) All() string {
	return consoleCollect.progText
}

// func TestMissingFileFromMenu(t *testing.T) {
// 	fmt.Println("aaa", Console, "aaa")
// 	Console.Note("hi")
// 	fmt.Println("bbb", Console.progText, "bbb")
// 	Console.Clear()
// 	fmt.Println("ccc", Console, "ccc")
// }
