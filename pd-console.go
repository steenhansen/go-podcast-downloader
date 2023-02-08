package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"

	"github.com/steenhansen/go-podcast-downloader-console/src/help"
)

func main() {
	diskSize, progBounds, cleanArgs := misc.InitProg(consts.MIN_DISK_BYTES)

	keyStream := make(chan string)

	// go func() {
	// 	fmt.Println("************* start sleep")
	// 	time.Sleep(time.Second * 31)
	// 	fmt.Println("************* stop sleep")
	// 	keyStream <- "a"
	// }()

	if len(cleanArgs) == 1 {
		fmt.Println(diskSize)
		for {
			report, err := menu.DisplayMenu(progBounds, keyStream, misc.KeyboardMenuChoice, rss.HttpMedia)
			if err != nil && !errors.Is(err, flaws.SStop) {
				panic(err)
			} else if report == "" {
				break // entered "Q" to quit
			}
			fmt.Println(report)
		}
	} else {
		arg1Lower := strings.ToLower(cleanArgs[1])
		if arg1Lower == consts.HELP_PLAIN || arg1Lower == consts.HELP_DASH || arg1Lower == consts.HELP_DASH_DASH {
			fmt.Println(help.HelpText())
		} else {
			fmt.Println(diskSize)
			report, err := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, rss.HttpMedia)
			if err != nil && !errors.Is(err, flaws.SStop) {
				panic(err)
			}
			fmt.Println(report, globals.Faults.All())
		}

	}
	fmt.Print("goodbye")
}
