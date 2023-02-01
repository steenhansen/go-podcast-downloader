package main

/*
go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
go run ./pd-console.go https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
*/

import (
	"errors"
	"fmt"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"

	"github.com/steenhansen/go-podcast-downloader-console/src/help"
)

func main() {
	diskSize, progBounds, cleanArgs := misc.InitProg(misc.DiskSpace, consts.MIN_DISK_BYTES)
	fmt.Println(diskSize)
	simKeyStream := make(chan string)

	// go func() {
	// 	fmt.Println("************* start sleep")
	// 	time.Sleep(time.Second * 31)
	// 	fmt.Println("************* stop sleep")
	// 	simKeyStream <- "a"
	// }()

	if len(cleanArgs) == 1 {
		for {
			report, err := menu.DisplayMenu(progBounds, simKeyStream, misc.KeyboardMenuChoice)
			if err != nil && !errors.Is(err, flaws.SStop) {
				panic(err)
			} else if report == "" {
				break // entered "Q" to quit
			}
			fmt.Println(report)

		}
	} else if strings.ToLower(cleanArgs[1]) == consts.HELP_ARG1 {
		fmt.Println(help.HelpText())
	} else {
		report, err := menu.ByNameOrUrl(cleanArgs, progBounds, simKeyStream)
		if err != nil && !errors.Is(err, flaws.SStop) {
			panic(err)
		}
		fmt.Println(report, globals.Faults.All())
	}

	fmt.Print("goodbye")
}
