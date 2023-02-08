package menu

import (
	"errors"

	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func ByNameOrUrl(cleanArgs []string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (podReport string, err error) {
	if feed.IsUrl(cleanArgs[1]) {
		feedUrl := cleanArgs[1]
		if len(cleanArgs) == 2 {
			podReport, err = terminal.AddByUrl(feedUrl, progBounds, keyStream, httpMedia) // go run ./ https://www.a.com/feed
		} else {
			podReport, err = terminal.AddByUrlAndName(feedUrl, cleanArgs, progBounds, keyStream, httpMedia) // go run ./ https://www.a.com/feed  My Fav Feed
		}
	} else {
		podReport, err = terminal.ReadByExistName(cleanArgs, progBounds, keyStream, httpMedia) // go run ./ My Fav Feed
	}
	return podReport, err
}

func DisplayMenu(progBounds models.ProgBounds, keyStream chan string, getMenuChoice models.ReadLineFn, httpMedia models.HttpFn) (string, error) {
	theMenu, _ := terminal.ShowNumberedChoices(progBounds)
	globals.Console.Note(theMenu)
	podReport, err := terminal.AfterMenu(progBounds, keyStream, getMenuChoice, httpMedia)
	if podReport == "" && err == nil {
		return "", nil // 'Q' entered to quit
	}
	if errors.Is(err, flaws.BadChoice) {
		return err.Error(), nil
	}
	if err != nil && !errors.Is(err, flaws.SStop) {
		return "", err
	}
	badFiles := globals.Faults.All()
	globals.Faults.Clear()
	allReport := podReport + "\n" + badFiles
	return allReport, err
}
