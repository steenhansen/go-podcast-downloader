package console

import (
	"errors"

	"podcast-downloader/src/dos/feed"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/terminal"

	"podcast-downloader/src/dos/flaws"
)

func ByNameOrUrl(cleanArgs []string, progBounds models.ProgBounds, keyStreamTest chan string, httpMedia models.HttpFn) (podReport string, podcastResults models.PodcastResults) {
	if feed.IsUrl(cleanArgs[1]) {
		feedUrl := cleanArgs[1]
		if len(cleanArgs) == 2 {
			podReport, podcastResults = terminal.AddByUrl(feedUrl, progBounds, keyStreamTest, httpMedia) // go run ./ https://www.a.com/feed
		} else {
			podReport, podcastResults = terminal.AddByUrlAndName(feedUrl, cleanArgs, progBounds, keyStreamTest, httpMedia) // go run ./ https://www.a.com/feed  My Fav Feed
		}
	} else {
		podReport, podcastResults = terminal.ReadByExistName(cleanArgs, progBounds, keyStreamTest, httpMedia) // go run ./ My Fav Feed
	}
	return podReport, podcastResults
}

func DisplayMenu(progBounds models.ProgBounds, keyStreamTest chan string, getMenuChoice models.ReadLineFn, httpMedia models.HttpFn) (string, bool, models.PodcastResults) {
	theMenu, _ := terminal.ShowNumberedChoices(progBounds)
	globals.Console.Note(theMenu)
	podReport, podcastResults := terminal.AfterMenu(progBounds, keyStreamTest, getMenuChoice, httpMedia)
	didQuit := false
	if podcastResults.WasCanceled && podcastResults.SeriousError == nil && podReport == "" {
		didQuit = true
	}
	return podReport, didQuit, podcastResults
}

func ShowResults(podReport string, podcastResults models.PodcastResults) {
	if podcastResults.SeriousError != nil {
		var flawError flaws.FlawError
		if errors.As(podcastResults.SeriousError, &flawError) {
			globals.Console.Note("\n" + flawError.Error() + "\n\n")
		} else {
			globals.Console.Note("\nUNKNOWN? " + podcastResults.SeriousError.Error() + "\n")
		}
	} else {

		if podcastResults.WasCanceled {
			globals.Console.Note("\nPodcast backup got canceled" + "\n")
		}
		globals.Console.Note(podReport + "\n")
		globals.Console.Note(globals.Faults.All() + "\n")

	}
}
