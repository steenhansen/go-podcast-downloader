package main

import (
	"fmt"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/stop"

	"github.com/steenhansen/go-podcast-downloader-console/src/help"
)

func main() {
	diskSize, progBounds, cleanArgs := misc.InitProg(consts.MIN_DISK_BYTES)
	keyStream := make(chan string)
	if len(cleanArgs) == 1 {
		fmt.Println(diskSize)
		for {
			podReport, didQuit, podcastResults := menu.DisplayMenu(progBounds, keyStream, stop.KeyboardMenuChoice, rss.HttpReal)
			//fmt.Println("podReport", podReport)
			//fmt.Println("didQuit", didQuit)
			//fmt.Println("podcastResults", podcastResults)
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
