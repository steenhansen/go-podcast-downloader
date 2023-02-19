package main

import (
	"fmt"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"

	"github.com/steenhansen/go-podcast-downloader-console/src/help"
)

func main() {
	diskSize, progBounds, cleanArgs := misc.InitProg(consts.MIN_DISK_BYTES)
	keyStream := make(chan string)
	if len(cleanArgs) == 1 {
		fmt.Println(diskSize)
		for {
			report, err := menu.DisplayMenu(progBounds, keyStream, misc.KeyboardMenuChoice, rss.HttpReal)
			if report == "" {
				break // entered "Q" to quit
			}
			menu.ShowResults(report, err)
		}
	} else {
		arg1Lower := strings.ToLower(cleanArgs[1])
		if arg1Lower == consts.HELP_PLAIN || arg1Lower == consts.HELP_DASH || arg1Lower == consts.HELP_DASH_DASH {
			fmt.Println(help.HelpText())
		} else {
			fmt.Println(diskSize)
			report, err := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, rss.HttpReal)
			menu.ShowResults(report, err)
		}

	}
	fmt.Print("goodbye")
}
