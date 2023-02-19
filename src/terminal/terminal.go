package terminal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
	"github.com/steenhansen/go-podcast-downloader-console/src/processes"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func ShowNumberedChoices(progBounds models.ProgBounds) (string, error) {
	podDirNames, thePodcasts, _, err := podcasts.AllPodcasts(progBounds.ProgPath)
	if err != nil {
		return "", err
	}
	if len(thePodcasts) == 0 {
		return noPodcastsAdded(flaws.FLAW_E_50)
	}
	podcastChoices, err := podcasts.PodChoices(progBounds.ProgPath, podDirNames)
	if err != nil {
		return "", err
	}
	theMenu := podcastChoices + " 'Q' or a number + enter: "
	return theMenu, nil
}

func AfterMenu(progBounds models.ProgBounds, keyStream chan string, getMenuChoice models.ReadLineFn, httpMedia models.HttpFn) (string, error) {
	podDirNames, thePodcasts, forceTitles, err := podcasts.AllPodcasts(progBounds.ProgPath)
	if err != nil {
		return "", err
	}
	podcastIndex, err := podcasts.ChoosePod(podDirNames, getMenuChoice)
	if podcastIndex == 0 && err == nil {
		return "", nil // 'Q' entered to quit
	}
	if err != nil {
		return "", err
	}
	globals.ForceTitle = forceTitles[podcastIndex-1] // derived from _origin-rss-url file
	mediaTitle := podDirNames[podcastIndex-1]
	podcastUrl := thePodcasts[podcastIndex-1]
	podcastResults := podcasts.DownloadPodcast(mediaTitle, podcastUrl, progBounds, keyStream, httpMedia)
	if podcastResults.Err != nil {
		return "", podcastResults.Err
	}
	podcastReport := doReport(podcastResults, string(podcastUrl), mediaTitle)
	return podcastReport, nil
}

func AddByUrl(podcastUrl string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (string, error) {
	rssXml, mediaTitles, rssFiles, rssSizes, err := podcasts.ReadRssUrl(podcastUrl, httpMedia, keyStream)
	if err != nil {
		return "", err
	}
	mediaTitle, err := rss.RssTitle(rssXml)
	if err != nil {
		return "", err
	}
	mediaPath, dirNotExist, err := media.InitFolder(progBounds.ProgPath, mediaTitle, podcastUrl)
	if err != nil {
		return "", err
	}
	if dirNotExist {
		globals.Console.Note("\nAdding '" + mediaTitle + "'\n\n")
	}
	podcastData := models.PodcastData{
		PodTitle:  mediaTitle,
		PodPath:   mediaPath,
		PodUrls:   rssFiles,
		PodSizes:  rssSizes,
		PodTitles: mediaTitles,
	}

	podcastReport, err := downloadReport(podcastUrl, podcastData, progBounds, keyStream, httpMedia)
	return podcastReport, err
}

func AddByUrlAndName(podcastUrl string, osArgs []string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (string, error) {
	_, mediaTitles, rssFiles, rssSizes, err := podcasts.ReadRssUrl(podcastUrl, httpMedia, keyStream)
	if err != nil {
		return "", err
	}
	mediaTitle := feed.PodcastName(osArgs)
	mediaPath, dirNotExist, err := media.InitFolder(progBounds.ProgPath, mediaTitle, podcastUrl)
	if err != nil {
		return "", err
	}
	if dirNotExist {
		globals.Console.Note("\nAdding '" + mediaTitle + "'\n\n")
	}
	podcastData := models.PodcastData{
		PodTitle:  mediaTitle,
		PodPath:   mediaPath,
		PodUrls:   rssFiles,
		PodSizes:  rssSizes,
		PodTitles: mediaTitles,
	}

	podcastReport, err := downloadReport(podcastUrl, podcastData, progBounds, keyStream, httpMedia)
	return podcastReport, err
}

func doReport(podcastResults models.PodcastResults, podcastUrl string, mediaTitle string) (podcastReport string) {
	savedFiles := podcastResults.SavedFiles
	varietyFiles := podcastResults.VarietyFiles
	podcastTime := podcastResults.PodcastTime
	var secRounded time.Duration
	if !misc.IsTesting(os.Args) {
		secRounded = podcastTime.Round(time.Second) // NB if testing all times are 0s
	}
	if savedFiles != 0 {
		addedNew := fmt.Sprintf("\nAdded %d new ", savedFiles)
		fileTypes := fmt.Sprintf("'%s' file(s) in %s \n", varietyFiles, secRounded)
		if len(varietyFiles) == 0 {
			fileTypes = fmt.Sprintf("files in %s \n", secRounded)
		}
		fromInto := fmt.Sprintf("From %s \nInto '%s' \n", podcastUrl, mediaTitle)
		podcastReport = addedNew + fileTypes + fromInto
	} else {
		podcastReport = "No changes detected"
	}
	return podcastReport
}

func downloadReport(url string, podcastData models.PodcastData, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (string, error) {
	podcastResults := processes.DownloadMedia(url, podcastData, progBounds, keyStream, httpMedia)
	if podcastResults.Err != nil && !errors.Is(podcastResults.Err, flaws.SKeyStop) {
		return "", podcastResults.Err
	}
	podcastReport := doReport(podcastResults, url, podcastData.PodTitle)
	return podcastReport, podcastResults.Err
}

func ReadByExistName(osArgs []string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) (string, error) {
	podcastTitle := feed.PodcastName(osArgs)
	mediaPath, mediaTitle, err := podcasts.FindPodcastDirName(progBounds.ProgPath, podcastTitle)
	if err != nil {
		return "", err
	}
	originRss := mediaPath + "/" + consts.URL_OF_RSS_FN
	urlBytes, err := os.ReadFile(originRss)
	if err != nil {
		return "", err
	}
	urlLines := string(urlBytes)
	urlStrings := misc.SplitByNewline(urlLines)
	urlStr := urlStrings[0]
	if len(urlStrings) > 1 {
		if urlStrings[1] == "--forceTitle" {
			globals.ForceTitle = true
		}
	}
	_, mediaTitles, mediaUrl, mediaSize, err := podcasts.ReadRssUrl(urlStr, httpMedia, keyStream) // _ == unused xml
	if err != nil {
		return "", err
	}
	podcastData := models.PodcastData{
		PodTitle:  mediaTitle,
		PodPath:   mediaPath,
		PodUrls:   mediaUrl,
		PodSizes:  mediaSize,
		PodTitles: mediaTitles,
	}
	podcastReport, err := downloadReport(urlStr, podcastData, progBounds, keyStream, httpMedia)
	return podcastReport, err
}
