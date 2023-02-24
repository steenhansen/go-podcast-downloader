package menu

import (
	"errors"

	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func ByNameOrUrl(cleanArgs []string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (podReport string, podcastResults models.PodcastResults) {
	if feed.IsUrl(cleanArgs[1]) {
		feedUrl := cleanArgs[1]
		if len(cleanArgs) == 2 {
			podReport, podcastResults = terminal.AddByUrl(feedUrl, progBounds, keyStream, httpMedia) // go run ./ https://www.a.com/feed
		} else {
			podReport, podcastResults = terminal.AddByUrlAndName(feedUrl, cleanArgs, progBounds, keyStream, httpMedia) // go run ./ https://www.a.com/feed  My Fav Feed
		}
	} else {
		podReport, podcastResults = terminal.ReadByExistName(cleanArgs, progBounds, keyStream, httpMedia) // go run ./ My Fav Feed
	}
	return podReport, podcastResults
}

func DisplayMenu(progBounds models.ProgBounds, keyStream chan string, getMenuChoice models.ReadLineFn, httpMedia models.HttpFn) (string, bool, models.PodcastResults) {
	theMenu, _ := terminal.ShowNumberedChoices(progBounds)
	globals.Console.Note(theMenu)
	podReport, podcastResults := terminal.AfterMenu(progBounds, keyStream, getMenuChoice, httpMedia)
	didQuit := false
	//	fmt.Println("DisplayMenu -- podReport", podReport == "")
	//fmt.Println("DisplayMenu -- podcastResults.SeriousError == nil", podcastResults.SeriousError == nil)
	//fmt.Println("DisplayMenu -- !podcastResults.WasCanceled", !podcastResults.WasCanceled)
	if podcastResults.WasCanceled && podcastResults.SeriousError == nil && podReport == "" {
		didQuit = true
	}
	return podReport, didQuit, podcastResults
}

func ShowResults(podReport string, podcastResults models.PodcastResults) {
	if podcastResults.SeriousError != nil {
		var flawError flaws.FlawError
		if errors.As(podcastResults.SeriousError, &flawError) {
			globals.Console.Note("\nSerious Error: " + flawError.Error() + "\n")
		} else {
			globals.Console.Note("\nSerious Error: UNKNOWN? " + podcastResults.SeriousError.Error() + "\n")
		}
	} else {

		if podcastResults.WasCanceled {
			globals.Console.Note("\nPodcast backup got canceled" + "\n")
		}
		globals.Console.Note(podReport + "\n")
		globals.Console.Note(globals.Faults.All() + "\n")

	}
}
