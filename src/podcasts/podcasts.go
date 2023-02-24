package podcasts

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/processes"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func ReadRssUrl(rssUrl string, httpMedia models.HttpFn) ([]byte, []string, []string, []int, error) {
	podcastXml, err := feed.ReadRss(rssUrl, httpMedia)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	podcastTitle, err := rss.RssTitle(podcastXml)
	if podcastTitle == "" {
		xmlStr := string(podcastXml[0:consts.FIRST_BYTES_OF_ERROR_PAGE])
		return badReadRssUrl(xmlStr, flaws.FLAW_E_62)
	} else if err != nil {
		return nil, nil, nil, nil, err
	}
	mediaTitles, mediaUrls, mediaSizes, err := rss.RssItems(podcastXml)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return podcastXml, mediaTitles, mediaUrls, mediaSizes, nil
}

func FindPodcastDirName(ProgPath, podcastTitle string) (string, string, error) {
	progDir, err := os.Open(ProgPath)
	if err != nil {
		return "", "", err
	}
	defer progDir.Close()
	dirOfPodcasts, err := progDir.Readdir(0)
	if err != nil {
		return "", "", err
	}
	lowerTitle := strings.ToLower(podcastTitle)
	for _, podDir := range dirOfPodcasts {
		if !podDir.Mode().IsRegular() {
			dirName := podDir.Name()
			lowerDir := strings.ToLower(dirName)
			if lowerDir == lowerTitle {
				mediaPath := ProgPath + "/" + podDir.Name()
				return mediaPath, dirName, nil
			}
		}
	}
	return badPodDirName(podcastTitle, flaws.FLAW_E_60)
}

func DownloadPodcast(mediaTitle, rssUrl string, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) models.PodcastResults {
	if feed.IsUrl(rssUrl) {
		_, mediaTitles, mediaUrls, mediaSizes, err := ReadRssUrl(rssUrl, httpMedia)
		if err != nil {
			return misc.EmptyPodcastResults(false, err)
		}
		mediaPath := progBounds.ProgPath + "/" + mediaTitle
		podcastData := models.PodcastData{
			PodTitle:  mediaTitle,
			PodPath:   mediaPath,
			PodUrls:   mediaUrls,
			PodSizes:  mediaSizes,
			PodTitles: mediaTitles,
		}
		podcastResults := processes.BackupPodcast(rssUrl, podcastData, progBounds, keyStream, httpMedia)
		return podcastResults
	}
	return badRssUrl(rssUrl, flaws.FLAW_E_61)
}

func PodChoices(ProgPath string, podDirNames []string) (podChoices string, err error) {
	var sizedStr string
	for podIndex, podcastDirName := range podDirNames {
		podCount, dirSize, err := countFiles(ProgPath, podcastDirName)
		if err != nil {
			return "", err
		}
		if dirSize < consts.GB_BYTES {
			mbs := dirSize / consts.MB_BYTES
			sizedStr = fmt.Sprintf("%dMB", mbs)
		} else {
			gbs := float64(dirSize) / float64(consts.GB_BYTES)
			sizedStr = fmt.Sprintf("%.2fGB", gbs)
		}
		podChoices += fmt.Sprintf("%2d |%4d files |%7s | %s\n", podIndex+1, podCount-1, sizedStr, podcastDirName)
	}
	return podChoices, nil
}

func ChoosePod(podDirNames []string, getMenuChoice models.ReadLineFn) (menuChoice int, err error) {
	lineInput := getMenuChoice()
	textChoice := strings.Trim(lineInput, "\r\n")
	textLower := strings.ToLower(textChoice)
	if textLower == consts.QUIT_KEY_LOWER {
		return 0, nil
	}
	menuChoice, _ = strconv.Atoi(textChoice)
	if menuChoice < 1 || menuChoice > len(podDirNames) {
		return badPodNumber(textChoice, flaws.FLAW_E_63)
	}
	return menuChoice, nil
}

func countFiles(progPath, dirName string) (fileCount int, dirSize int64, err error) {
	dirPath := progPath + "/" + dirName
	dirFiles, err := misc.FilesInDir(dirPath)
	if err != nil {
		return 0, 0, err
	}
	for _, mediaFile := range dirFiles {
		if mediaFile.Mode().IsRegular() {
			dirSize = dirSize + mediaFile.Size()
			fileCount++
		}
	}
	return fileCount, dirSize, nil
}

func AllPodcasts(progPath string) ([]string, []string, []bool, error) {
	progDir, err := os.Open(progPath)
	if err != nil {
		return nil, nil, nil, err
	}
	defer progDir.Close()
	podcastDirs, err := progDir.Readdir(0)
	if err != nil {
		return nil, nil, nil, err
	}
	podDirNames := make([]string, 0)
	allFeeds := make([]string, 0)
	forceTitles := make([]bool, 0)
	for _, dir := range podcastDirs {
		if !dir.Mode().IsRegular() {
			dirName := dir.Name()
			if dirName != consts.SOURCE_FOLDER {
				rssPath := progPath + "/" + dir.Name() + "/" + consts.URL_OF_RSS_FN
				urlBytes, err := os.ReadFile(rssPath)
				urlLines := string(urlBytes)
				urlStrings := misc.SplitByNewline(urlLines)
				urlStr := urlStrings[0]
				if len(urlStrings) > 1 && urlStrings[1] == "--forceTitle" {
					forceTitles = append(forceTitles, true)
				} else {
					forceTitles = append(forceTitles, false)
				}
				if err == nil {
					podDirNames = append(podDirNames, dirName)
					allFeeds = append(allFeeds, urlStr)
				}

			}
		}
	}
	return podDirNames, allFeeds, forceTitles, nil
}
