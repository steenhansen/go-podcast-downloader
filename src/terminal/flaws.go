package terminal

import (
	"errors"

	"github.com/steenhansen/go-podcast-downloader/src/flaws"
)

func noPodcastsAdded(locId string) (string, error) {
	locStat := errors.New(locId)
	err := flaws.EmptyPodcasts.MakeFlaw("add some podcasts feeds first")
	multiErr := errors.Join(locStat, err)
	return "", multiErr
}
