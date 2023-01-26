package menu

import (
	"errors"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func AddFeed(cleanArgs []string, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	var report string
	var err error
	if feed.IsUrl(cleanArgs[1]) {
		feedUrl := cleanArgs[1]
		if len(cleanArgs) == 2 {
			report, err = terminal.AddByUrl(feedUrl, progBounds, simKeyStream, mediaFix) // go run ./ https://www.a.com/feed
		} else {
			report, err = terminal.AddByUrlAndName(feedUrl, cleanArgs, progBounds, simKeyStream, mediaFix) // go run ./ https://www.a.com/feed  My Fav Feed
		}
	} else {

		report, err = terminal.ReadByExistName(cleanArgs, progBounds, simKeyStream, mediaFix) // go run ./ My Fav Feed
	}
	return report, err
}

// func continueError(err error) bool {
// 	return errors.Is(err, flaws.NoPodcasts) ||
// 		errors.Is(err, flaws.BadChoice) ||
// 		errors.Is(err, flaws.SStop) ||
// 		errors.Is(err, flaws.BadUrl) ||
// 		errors.Is(err, flaws.BadContent)
// }

func DisplayMenu(progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	report, err := terminal.ShowNumberedChoices(progBounds, simKeyStream, mediaFix)
	if report == "" && err == nil {
		return "", nil // 'Q' entered to quit
	}
	if err != nil {
		if errors.Is(err, flaws.LowDisk) {
			return "", err
		}
		if errors.Is(err, flaws.BadChoice) {
			return err.Error(), nil
		}
		return "", err
	}
	badFiles := ""
	for _, mediaError := range mediaFix {
		badFiles = badFiles + "\t" + mediaError.Error() + "\n"
	}
	report = report + "\n" + badFiles
	return report, nil
}
