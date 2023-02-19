package podcasts

import (
	"errors"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
)

func badPodDirName(podcastTitle, locId string) (string, string, error) {
	locStat := errors.New(locId)
	err := flaws.InvalidPodcastName.MakeFlaw(podcastTitle)
	multiErr := errors.Join(locStat, err)
	return "", "", multiErr
}

func badRssUrl(rssUrl, locId string) models.PodcastResults {
	locStat := errors.New(locId)
	err := flaws.InvalidPodcastName.MakeFlaw(rssUrl)
	multiErr := errors.Join(locStat, err)
	emptyPodcastResults := misc.EmptyPodcastResults(multiErr)
	return emptyPodcastResults
}

func badReadRssUrl(xmlStr, locId string) ([]byte, []string, []string, []int, error) {
	locStat := errors.New(locId)
	err := flaws.InvalidXML.MakeFlaw(xmlStr)
	multiErr := errors.Join(locStat, err)
	return nil, nil, nil, nil, multiErr
}

func badPodNumber(textChoice, locId string) (int, error) {
	locStat := errors.New(locId)
	err := flaws.BadChoice.MakeFlaw(textChoice)
	multiErr := errors.Join(locStat, err)
	return 0, multiErr
}
