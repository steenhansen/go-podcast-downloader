package main

import (
	"fmt"
	"strings"

	"podcast-downloader/src/dos/console"
	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/initialize"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/stop"

	"podcast-downloader/src/dos/help"
)

func main() { // console version
	initialize.AddNasa()
	diskSize, progBounds, cleanArgs := misc.InitProg()
	keyStreamTest := make(chan string)
	if len(cleanArgs) == 1 {
		fmt.Println(diskSize)

		for {
			podReport, didQuit, podcastResults := console.DisplayMenu(progBounds, keyStreamTest, stop.KeyboardMenuChoice, rss.HttpReal)
			if didQuit {
				break
			}
			console.ShowResults(podReport, podcastResults)
		}
	} else {
		arg1Lower := strings.ToLower(cleanArgs[1])
		if arg1Lower == consts.HELP_PLAIN || arg1Lower == consts.HELP_DASH || arg1Lower == consts.HELP_DASH_DASH {
			fmt.Println(help.HelpText())
		} else {
			fmt.Println(diskSize)
			podReport, podcastResults := console.ByNameOrUrl(cleanArgs, progBounds, keyStreamTest, rss.HttpReal)
			console.ShowResults(podReport, podcastResults)
		}

	}
	globals.Console.Note(consts.GOOD_BYE_MESS)
}
