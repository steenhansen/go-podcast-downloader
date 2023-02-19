package globals

import (
	"fmt"
)

type ConsoleCollect struct {
	progText string
}

const CLEAR_SCREEN = "\033[H\033[2J"
const CLEAR_LINE = "\033[2K\r"

// Catch console output so tests can verify
var Console = ConsoleCollect{progText: ""}

func (consoleCollect *ConsoleCollect) Note(progressStr string) {
	fmt.Print(CLEAR_LINE + progressStr)
	if progressStr != CLEAR_SCREEN {
		consoleCollect.progText = consoleCollect.progText + progressStr
	}
}

func (consoleCollect *ConsoleCollect) Clear() {
	consoleCollect.progText = ""
}

func (consoleCollect *ConsoleCollect) All() string {
	return consoleCollect.progText
}
