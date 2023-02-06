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
	"github.com/steenhansen/go-podcast-downloader-console/src/processes"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/varieties"
)

func ReadRssUrl(rssUrl string, httpMedia consts.HttpFunc) ([]byte, []string, []int, error) {
	podcastXml, err := feed.ReadRss(rssUrl, httpMedia)
	if err != nil {
		return nil, nil, nil, err
	}
	podcastTitle, err := rss.RssTitle(podcastXml)
	if podcastTitle == "" {
		xmlStr := string(podcastXml[0:100])
		return nil, nil, nil, flaws.InvalidXML.StartError(xmlStr)
	} else if err != nil {
		return nil, nil, nil, err
	}
	mediaUrls, mediaSizes, err := rss.RssItems(podcastXml)
	if err != nil {
		return nil, nil, nil, err
	}
	return podcastXml, mediaUrls, mediaSizes, nil
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
	return "", "", flaws.NoMatchName.StartError(podcastTitle)
}

func DownloadPodcast(mediaTitle, rssUrl string, progBounds consts.ProgBounds, keyStream chan string, httpMedia consts.HttpFunc) consts.PodcastResults {
	if feed.IsUrl(rssUrl) {

		_, mediaUrls, mediaSizes, err := ReadRssUrl(rssUrl, httpMedia)
		if err != nil {
			return misc.EmptyPodcastResults(err)
		}
		mediaPath := progBounds.ProgPath + "/" + mediaTitle
		podcastData := consts.PodcastData{
			PodTitle: mediaTitle,
			PodPath:  mediaPath,
			PodUrls:  mediaUrls,
			PodSizes: mediaSizes,
		}
		podcastResults := processes.DownloadMedia(rssUrl, podcastData, progBounds, keyStream, httpMedia)
		return podcastResults
	}
	return misc.EmptyPodcastResults(flaws.InvalidRssURL.StartError(rssUrl))
}

func PodChoices(ProgPath string, podDirNames []string) (podChoices string, err error) {
	var sizedStr string
	for podIndex, podcastDirName := range podDirNames {
		podCount, dirSize, fileTypes, err := countFiles(ProgPath, podcastDirName)
		if fileTypes == "" && strings.Contains(podcastDirName, "[") {
			nameParts := strings.Split(podcastDirName, "[")
			extensionInName := nameParts[1]
			fileTypes = extensionInName[0:3]
		}
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
		podChoices += fmt.Sprintf("%2d | %16s |%4d files |%7s | %s\n", podIndex+1, fileTypes, podCount-1, sizedStr, podcastDirName)
	}
	return podChoices, nil
}

func ChoosePod(podDirNames []string, getMenuChoice consts.ReadLineFunc) (menuChoice int, err error) {
	lineInput := getMenuChoice()
	textChoice := strings.Trim(lineInput, "\r\n")
	textLower := strings.ToLower(textChoice)
	if textLower == consts.QUIT_KEY_LOWER {
		return 0, nil
	}
	menuChoice, _ = strconv.Atoi(textChoice)
	if menuChoice < 1 || menuChoice > len(podDirNames) {
		return 0, flaws.BadChoice.StartError(textChoice)
	}
	return menuChoice, nil
}

func countFiles(progPath, dirName string) (fileCount int, dirSize int64, varietyFiles string, err error) {
	varietySet := varieties.VarietiesSet{}
	dirPath := progPath + "/" + dirName
	podDir, err := os.Open(dirPath)
	if err != nil {
		return 0, 0, "", err
	}
	defer podDir.Close()

	dirFiles, err := podDir.Readdir(0)
	if err != nil {
		return 0, 0, "", err
	}

	for _, mediaFile := range dirFiles {
		if mediaFile.Mode().IsRegular() {
			varietySet.AddVariety(mediaFile.Name())
			dirSize = dirSize + mediaFile.Size()
			fileCount++
		}
	}
	varietyFiles = varietySet.VarietiesString(" ")
	return fileCount, dirSize, varietyFiles, nil
}

func AllPodcasts(progPath string) ([]string, []string, error) {
	progDir, err := os.Open(progPath)
	if err != nil {
		return nil, nil, err
	}
	defer progDir.Close()
	podcastDirs, err := progDir.Readdir(0)
	if err != nil {
		return nil, nil, err
	}
	podDirNames := make([]string, 0)
	allFeeds := make([]string, 0)
	for _, dir := range podcastDirs {
		if !dir.Mode().IsRegular() {
			dirName := dir.Name()
			if dirName != consts.SOURCE_FOLDER {
				rssPath := progPath + "/" + dir.Name() + "/" + consts.URL_OF_RSS_FN
				rssUrl, err := os.ReadFile(rssPath)
				if err == nil {
					podDirNames = append(podDirNames, dirName)
					allFeeds = append(allFeeds, string(rssUrl))
				}

			}
		}
	}
	return podDirNames, allFeeds, nil
}
