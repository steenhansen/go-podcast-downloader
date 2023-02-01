package menu

import (
	"errors"
	"fmt"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func ByNameOrUrl(cleanArgs []string, progBounds consts.ProgBounds, simKeyStream chan string) (podReport string, err error) {
	if feed.IsUrl(cleanArgs[1]) {
		feedUrl := cleanArgs[1]
		if len(cleanArgs) == 2 {
			podReport, err = terminal.AddByUrl(feedUrl, progBounds, simKeyStream) // go run ./ https://www.a.com/feed
		} else {
			podReport, err = terminal.AddByUrlAndName(feedUrl, cleanArgs, progBounds, simKeyStream) // go run ./ https://www.a.com/feed  My Fav Feed
		}
	} else {
		podReport, err = terminal.ReadByExistName(cleanArgs, progBounds, simKeyStream) // go run ./ My Fav Feed
	}
	return podReport, err
}

func DisplayMenu(progBounds consts.ProgBounds, simKeyStream chan string, getMenuChoice consts.ReadLineFunc) (string, error) {
	theMenu, _ := terminal.ShowNumberedChoices(progBounds)
	fmt.Print(theMenu)
	podReport, err := terminal.AfterMenu(progBounds, simKeyStream, getMenuChoice)
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
