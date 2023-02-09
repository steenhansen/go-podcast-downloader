package media

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/varieties"
)

func dirTitle(podTitle, rssUrl string) string {
	if podTitle == "" {
		podTitle = rssUrl
	}
	xmlEscaped := regexp.MustCompile(`&[^;]*;`) // &amp;
	safeXml := xmlEscaped.ReplaceAllLiteralString(podTitle, "")

	multSpaces := regexp.MustCompile(`\s\s+`)
	singleSpaces := multSpaces.ReplaceAllLiteralString(safeXml, " ")

	alphaRe := regexp.MustCompile(`[^a-zA-Z0-9 _-]+`)
	safeTitle := alphaRe.ReplaceAllLiteralString(singleSpaces, "")
	return safeTitle
}

func InitFolder(progPath, podTitle, rssUrl string) (string, bool, error) {
	safeTitle := dirTitle(podTitle, rssUrl)
	containDir := progPath + "/" + safeTitle
	dirNotExist := false
	if _, err := os.Stat(containDir); os.IsNotExist(err) {
		if err := os.Mkdir(containDir, os.ModePerm); err != nil {
			return "", false, flaws.CantCreateDirSerious.ContinueError(containDir, err)
		}
		dirNotExist = true
	}
	originRss := containDir + "/" + consts.URL_OF_RSS_FN
	rssAddrFile, err := os.Create(originRss)
	if err != nil {
		return "", false, flaws.CantCreateFileSerious.ContinueError(originRss, err)
	}
	defer rssAddrFile.Close()
	_, err = rssAddrFile.Write([]byte(rssUrl))
	if err != nil {
		return "", false, flaws.CantWriteFileSerious.ContinueError(originRss, err)
	}
	return containDir, dirNotExist, nil
}

func chooseName(finalFileName string, podPath string, podTitle string, podcastData models.PodcastData) (filePath string) {
	filePath = podPath + "/"
	if podTitle == "" {
		filePath += finalFileName // none found yet
	} else {
		mediaUrlsSet := make(map[string]string)
		for _, podUrl := range podcastData.PodUrls {
			fileName := rss.NameOfFile(podUrl)
			mediaUrlsSet[fileName] = fileName
		}
		// NHK Japan always has 1 mp3
		if len(podcastData.PodUrls) > 1 && len(mediaUrlsSet) == 1 { // sysk.com/redirect.mp3
			fileExt := varieties.FindVariety(finalFileName)
			var re = regexp.MustCompile(consts.BAD_FILE_CHAR_AND_DOT)
			goodChars := re.ReplaceAllString(podTitle, "")
			trimmedName := strings.TrimSpace(goodChars)
			filePath += trimmedName + "." + fileExt
		} else {
			filePath += finalFileName // nearly every podcast
		}

	}
	return filePath
}

func SaveDownloadedMedia(ctx context.Context, podcastData models.PodcastData, mediaStream chan<- models.MediaEnclosure, limitFlag int, httpMedia models.HttpFn) (int, string, error) {
	varietySet := varieties.VarietiesSet{}
	possibleFiles := 0
limitCancel:
	for mediaIndex, mediaUrl := range podcastData.PodUrls {
		select {
		case <-ctx.Done():
			break limitCancel
		default:
			possibleFiles++
			finalFileName, err := rss.FinalMediaName(ctx, mediaUrl, httpMedia)
			if err != nil {
				return 0, "", err
			}
			podTitle := ""
			if len(podcastData.PodTitles) > mediaIndex {
				podTitle = podcastData.PodTitles[mediaIndex]
			}

			filePath := chooseName(finalFileName, podcastData.PodPath, podTitle, podcastData)
			varietySet.AddVariety(finalFileName)
			if _, err = os.Stat(filePath); err != nil {
				if os.IsNotExist(err) {
					newMedia := models.MediaEnclosure{
						EnclosureUrl:  mediaUrl,
						EnclosurePath: filePath,
						EnclosureSize: podcastData.PodSizes[mediaIndex],
					}
					mediaStream <- newMedia
					limitFlag--
					if limitFlag == 0 {
						break limitCancel
					}
				} else {
					return 0, "", err
				}
			}
		}
	}
	varietyFiles := varietySet.VarietiesString(" ")
	return possibleFiles, varietyFiles, nil
}
