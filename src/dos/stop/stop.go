package stop

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/feed"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"

	"github.com/eiannone/keyboard"
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

func Go_seriousError(ctx context.Context, cancel context.CancelFunc,
	errorStream <-chan models.MediaError,
	seriousStream chan<- error, signalEndSerious <-chan bool, downloadEpisodeErrorEvent func(string)) {
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
			globals.Console.Note("ERROR: " + err.Error())
			globals.Console.Note(feed.ShowError(fileName))
			//	values.An_exe_debug_error_message = "exe_debug Go_seriousError " + fileName + " - " + err.Error()\
			if !errors.Is(err, flaws.NoGuiKeyboard) {
				downloadEpisodeErrorEvent(fileName)
			}
		}
	}
	misc.ChannelLog("\t\t\t Go_seriousError END")
}

func Go_ctxDone(ctx context.Context) {
	misc.ChannelLog("\t Go_ctxDone START")
	<-ctx.Done()
	misc.ChannelLog("\t Go_ctxDone END")
}

// keyStreamTest <-chan string is so that a test can simulate stopping
func Go_stopKey(cancel context.CancelFunc, KeyEventsReal <-chan keyboard.KeyEvent, keyStreamTest <-chan string, signalEndStop <-chan bool, afterDownloadEpisodeEvent func(string)) {
	misc.ChannelLog("\t\t Go_stopKey START")
	stoppedByKey := false
	guiTitleCount := 0
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
				misc.ChannelLog("\t Go_stopKey wasStopKeyed <- true")
				globals.Console.Note("stopping ...    \n")
				cancel()
				break keyboardEnd
			}
		case simKey := <-keyStreamTest:
			stoppedByKey = true
			misc.ChannelLog("\t\t TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			globals.Console.Note("TESTING - downloading stopped by simulated key press of '" + simKey + "'\n")
			cancel()
			break keyboardEnd
		default:
			spinBusy(guiTitleCount, afterDownloadEpisodeEvent)
			guiTitleCount++

		}
	}
	if stoppedByKey {
		misc.ChannelLog("\t\t Go_stopKey already STOPPED waiting for <-signalEndStop")
		<-signalEndStop
		misc.ChannelLog("\t\t Go_stopKey STOPPED finished")
	}
	misc.ChannelLog("\t\t Go_stopKey END")
}

const MASK_0_TO_31 = 31

func spinBusy(guiTitleCount int, afterDownloadEpisodeEvent func(string)) {
	numberGuiChars := guiTitleCount & MASK_0_TO_31
	if !consts.IsTesting(os.Args) {
		timeNow := time.Now()
		nanoNow := timeNow.Nanosecond()
		zeroTo99 := nanoNow / 10_000_000
		if zeroTo99 < 25 {
			fmt.Print(SPIN_0)
			title0Repeats := strings.Repeat(SPIN_0, numberGuiChars)
			afterDownloadEpisodeEvent(title0Repeats)
		} else if zeroTo99 < 50 {
			fmt.Print(SPIN_1)
			title1Repeats := strings.Repeat(SPIN_1, numberGuiChars)
			afterDownloadEpisodeEvent(title1Repeats)
		} else if zeroTo99 < 75 {
			fmt.Print(SPIN_2)
			title2Repeats := strings.Repeat(SPIN_2, numberGuiChars)
			afterDownloadEpisodeEvent(title2Repeats)
		} else {
			fmt.Print(SPIN_3)
			title3Repeats := strings.Repeat(SPIN_3, numberGuiChars)
			afterDownloadEpisodeEvent(title3Repeats)
		}
	}
}
