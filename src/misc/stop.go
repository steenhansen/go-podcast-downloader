package misc

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
)

const SPIN_0 = "\r|"
const SPIN_1 = "\r/"
const SPIN_2 = "\r-"
const SPIN_3 = "\r\\"

const SPIN_TYPES = 4

func KeyboardMenuChoice() string {
	keyboardReader := bufio.NewReader(os.Stdin)
	inputText, _ := keyboardReader.ReadString('\n')
	return inputText
}

// keyStreamTest <-chan string is so that a test can simulate stopping
func GoStopKey(ctx context.Context, cancel context.CancelFunc, KeyEventsReal <-chan keyboard.KeyEvent, keyStreamTest <-chan string) {
keyboardCancel:
	for {
		select {
		case event := <-KeyEventsReal:
			keyChar := string(event.Rune)
			keyLower := strings.ToLower(keyChar)
			if keyLower == consts.STOP_KEY_LOWER {
				globals.Console.Note("stopping ... \n")
				globals.StopingOnSKey = true
				cancel()
				break keyboardCancel
			}
		case simKey := <-keyStreamTest:
			globals.Console.Note("TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			cancel()
			break keyboardCancel
		case <-ctx.Done():
			// has timedout from RSS_MAX_READ_FILE_TIME or MEDIA_MAX_READ_FILE_TIME
			cancel()
			break keyboardCancel
		default:
			randInt := rand.Intn(SPIN_TYPES)
			if randInt == 0 {
				fmt.Print(SPIN_0)
			} else if randInt == 1 {
				fmt.Print(SPIN_1)
			} else if randInt == 2 {
				fmt.Print(SPIN_2)
			} else {
				fmt.Print(SPIN_3)
			}
		}

	}
}
