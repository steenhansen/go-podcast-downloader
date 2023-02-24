package stop

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
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

func Go_seriousError(ctx context.Context, cancel context.CancelFunc, errorStream <-chan models.MediaError, seriousStream chan<- error, signalEndSerious <-chan bool) {
	misc.ChannelLog("\t\t\t Go_seriousError START")
seriousEnd:
	for {
		select {
		case <-signalEndSerious:
			misc.ChannelLog("\t\t\t Go_seriousError <-signalEndSerious")
			break seriousEnd
		case mediaError := <-errorStream:
			err := mediaError.OrgErr
			if flaws.IsSerious(err) { // don't crash on a missing media file
				seriousStream <- err
				cancel()
			}
			fileName := rss.NameOfFile(mediaError.EnclosurePath)
			globals.Faults.Note(mediaError.EnclosurePath, err)
			globals.Console.Note(feed.ShowError(fileName))
		}
	}
	misc.ChannelLog("\t\t\t Go_seriousError END")
}

func Go_ctxDone(ctx context.Context, signalEndDerive chan<- bool, wasStopKeyed <-chan bool) {
	misc.ChannelLog("\t Go_ctxDone START")
	<-ctx.Done()
	select {
	case <-wasStopKeyed:
		misc.ChannelLog("\t\t\t\t Go_deriveFilenames end caused by ctx")
		misc.ChannelLog("\t Go_ctxDone wasStopKeyed END")
	default:
		misc.ChannelLog("\t Go_ctxDone signalEndDerive END")
		signalEndDerive <- true
	}
}

// keyStreamTest <-chan string is so that a test can simulate stopping
func Go_stopKey(cancel context.CancelFunc, KeyEventsReal <-chan keyboard.KeyEvent,
	keyStreamTest <-chan string, signalEndStop <-chan bool, wasStopKeyed chan<- bool) {
	misc.ChannelLog("\t\t Go_stopKey START")
	stoppedByKey := false
keyboardEnd:
	for {
		select {
		case <-signalEndStop:
			misc.ChannelLog("\t Go_stopKey NOT STOPPED <-signalEndStop")
			break keyboardEnd
		case event := <-KeyEventsReal:
			keyChar := string(event.Rune)
			keyLower := strings.ToLower(keyChar)
			if keyLower == consts.STOP_KEY_LOWER {
				stoppedByKey = true
				misc.ChannelLog("\t Go_stopKey wasStopKeyed")
				wasStopKeyed <- true
				misc.ChannelLog("\t Go_stopKey wasStopKeyed <- true")
				globals.Console.Note("stopping ...    \n")
				cancel()
				break keyboardEnd
			}
		case simKey := <-keyStreamTest:
			stoppedByKey = true
			wasStopKeyed <- true
			globals.Console.Note("TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			cancel()
			break keyboardEnd
		default:
			spinBusy()
		}
	}
	if stoppedByKey {
		misc.ChannelLog("\t\t Go_stopKey already STOPPED waiting for <-signalEndStop")
		<-signalEndStop
		misc.ChannelLog("\t\t Go_stopKey STOPPED finished")
	}
	misc.ChannelLog("\t\t Go_stopKey END")
}

func spinBusy() {
	if !consts.IsTesting(os.Args) {
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
