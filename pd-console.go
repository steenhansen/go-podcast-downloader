package main

import (
	"fmt"
	"strings"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/stop"

	"github.com/steenhansen/go-podcast-downloader/src/help"
)

func main() {
	diskSize, progBounds, cleanArgs := misc.InitProg(consts.MIN_DISK_BYTES)
	keyStream := make(chan string)
	if len(cleanArgs) == 1 {
		fmt.Println(diskSize)
		for {
			podReport, didQuit, podcastResults := menu.DisplayMenu(progBounds, keyStream, stop.KeyboardMenuChoice, rss.HttpReal)
			if didQuit {
				break
			}
			menu.ShowResults(podReport, podcastResults)
		}
	} else {
		arg1Lower := strings.ToLower(cleanArgs[1])
		if arg1Lower == consts.HELP_PLAIN || arg1Lower == consts.HELP_DASH || arg1Lower == consts.HELP_DASH_DASH {
			fmt.Println(help.HelpText())
		} else {
			fmt.Println(diskSize)
			podReport, podcastResults := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, rss.HttpReal)
			menu.ShowResults(podReport, podcastResults)
		}

	}
	globals.Console.Note(consts.GOOD_BYE_MESS)
}
