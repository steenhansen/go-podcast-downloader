package globals

import (
	"fmt"
	"sync"
)

type ConsoleCollect struct {
	progText string
}

var consoleMutex = sync.RWMutex{}

const CLEAR_SCREEN = "\033[H\033[2J"
const CLEAR_LINE = "\033[2K\r"

// Catch console output so tests can verify
var Console = ConsoleCollect{progText: ""}

// n.b. must print a "\n" at end of line or spinBusy() will overwrite
func (consoleCollect *ConsoleCollect) Note(progressStr string) {
	consoleMutex.Lock()
	fmt.Print(CLEAR_LINE + progressStr)
	if progressStr != CLEAR_SCREEN {
		consoleCollect.progText = consoleCollect.progText + progressStr
	}
	consoleMutex.Unlock()
}

func (consoleCollect *ConsoleCollect) Clear() {
	consoleCollect.progText = ""
}

func (consoleCollect *ConsoleCollect) All() string {
	return consoleCollect.progText
}
