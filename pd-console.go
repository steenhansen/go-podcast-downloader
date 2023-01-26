package main

/*
go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
go run ./pd-console.go https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
*/

//https://stackoverflow.com/questions/61845013/package-xxx-is-not-in-goroot-when-building-a-go-project
import (
	"fmt"
	"os"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"

	"github.com/steenhansen/go-podcast-downloader-console/src/help"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
)

func main() {
	free, size, percent := misc.DiskSpace()
	fmt.Printf("Current disk has %s free from a total %s which is %s full\n", free, size, percent)
	raceArgs, err := misc.DelRace(os.Args)
	if err != nil {
		panic(err)
	}

	limitFlag, tempArgs, err := misc.LimitArg(raceArgs)
	if err != nil {
		panic(err)
	}
	loadFlag, cleanArgs, err := misc.LoadArg(tempArgs)
	if err != nil {
		panic(err)
	}

	path := media.CurDir()
	progBounds := consts.ProgBounds{
		ProgPath:    path,
		LoadOption:  loadFlag,
		LimitOption: limitFlag,
		MinDisk:     consts.MIN_DISK_BYTES,
		//MinDisk:   consts.MIN_DISK_FAIL_BYTES,
	}
	simKeyStream := make(chan string)
	mediaFix := map[string]error{}
	// go func() {
	// 	fmt.Println("************* start sleep")
	// 	time.Sleep(time.Second * 31)
	// 	fmt.Println("************* stop sleep")
	// 	simKeyStream <- "a"
	// }()

	if len(cleanArgs) == 1 {
		for {
			report, err := menu.DisplayMenu(progBounds, simKeyStream, mediaFix)
			if err != nil {
				panic(err)
			}
			if report == "" {
				break // entered "Q" to quit
			}
			fmt.Println(report)
		}
	} else if strings.ToLower(cleanArgs[1]) == consts.HELP_ARG1 {
		helpText := help.HelpText()
		fmt.Println(helpText)
	} else {
		report, err := menu.AddFeed(cleanArgs, progBounds, simKeyStream, mediaFix)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(report)
		}
	}

	fmt.Println("THE ERRORS: ", mediaFix)
	fmt.Print("goodbye")
}
