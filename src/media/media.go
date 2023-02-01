package media

import (
	"context"
	"os"
	"regexp"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
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

func SaveDownloadedMedia(ctx context.Context, podcastData consts.PodcastData, mediaStream chan<- consts.MediaEnclosure, limitFlag int) (int, string, error) {
	varietySet := varieties.VarietiesSet{}
	possibleFiles := 0
limitCancel:
	for mediaIndex, mediaUrl := range podcastData.PodUrls {
		select {
		case <-ctx.Done():
			break limitCancel
		default:
			possibleFiles++
			finalFileName, err := rss.FinalMediaName(ctx, mediaUrl)
			if err != nil {
				return 0, "", err
			}
			filePath := podcastData.PodPath + "/" + finalFileName
			_, err = os.Stat(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					newMedia := consts.MediaEnclosure{
						EnclosureUrl:  mediaUrl,
						EnclosurePath: filePath,
						EnclosureSize: podcastData.PodSizes[mediaIndex],
					}
					mediaStream <- newMedia
					varietySet.AddVariety(finalFileName)
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
