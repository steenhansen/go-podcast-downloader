package media

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
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
			return doesntExist(containDir, flaws.FLAW_E_70, err)
		}
		dirNotExist = true
	}
	originRss := containDir + "/" + consts.URL_OF_RSS_FN
	rssAddrFile, err := os.Create(originRss)
	if err != nil {
		return cannotCreate(originRss, flaws.FLAW_E_71, err)
	}
	defer rssAddrFile.Close()
	_, err = rssAddrFile.Write([]byte(rssUrl)) // media files change size on ad injection
	if err != nil {
		return writeError(originRss, flaws.FLAW_E_72, err)
	}
	if globals.ForceTitle {
		rssAddrFile.Write([]byte("\n--forceTitle"))
	}
	return containDir, dirNotExist, nil
}

func chooseName(finalFileName string, podPath string, podTitle string, podcastData models.PodcastData) (filePath string) {
	filePath = podPath + "/"
	filePieces := strings.Split(finalFileName, ".")
	fileExt := filePieces[len(filePieces)-1]
	noSlashes := strings.ReplaceAll(podTitle, "/", "-") // for dates
	invFNameChars := regexp.MustCompile(consts.BAD_FILE_CHAR_AND_DOT)
	noInvFNameChars := invFNameChars.ReplaceAllString(noSlashes, "")
	multSpaces := regexp.MustCompile(`\s\s+`)
	goodSpaces := multSpaces.ReplaceAllString(noInvFNameChars, " ")
	shortName := goodSpaces
	if len(goodSpaces) > consts.MAX_TITLE_LEN {
		shortName = goodSpaces[:consts.MAX_TITLE_LEN]
	}
	trimmedName := strings.TrimSpace(shortName)
	if podTitle == "" {
		filePath += finalFileName
	} else if globals.ForceTitle {
		filePath += trimmedName + "." + fileExt
	} else {
		mediaUrlsSet := make(map[string]string)
		for _, podUrl := range podcastData.PodUrls {
			fileName := rss.NameOfFile(podUrl)
			mediaUrlsSet[fileName] = fileName
		}
		// NHK Japan always has 1 mp3
		if len(podcastData.PodUrls) > 1 && len(mediaUrlsSet) == 1 { // sysk.com/redirect.mp3
			filePath += trimmedName + "." + fileExt
		} else {
			filePath += finalFileName // nearly every podcast
		}
	}
	return filePath
}

func Go_deriveFilenames(ctx context.Context, podcastData models.PodcastData, mediaStream chan<- models.MediaEnclosure,
	limitFlag int, httpMedia models.HttpFn) (int, error) {
	misc.ChannelLog("\t\t\t Go_deriveFilenames START")
	possibleFiles := 0
	haveCount := 0
limitCancel:
	for mediaIndex, mediaUrl := range podcastData.PodUrls {
		possibleFiles++
		finalFileName, err := rss.FinalMediaName(ctx, mediaUrl, httpMedia)
		misc.ChannelLog("\t\t\t\t\t\t Go_deriveFilenames " + strconv.Itoa(possibleFiles) + " " + finalFileName)
		if err != nil {
			return 0, err
		}
		podTitle := ""
		if len(podcastData.PodTitles) > mediaIndex {
			podTitle = podcastData.PodTitles[mediaIndex]
		}
		filePath := chooseName(finalFileName, podcastData.PodPath, podTitle, podcastData)
		nameOfFile := rss.NameOfFile(filePath)
		misc.ChannelLog("\t\t\t\t\t\t Go_deriveFilenames " + strconv.Itoa(possibleFiles) + " " + nameOfFile)
		if _, err = os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				newMedia := models.MediaEnclosure{
					EnclosureUrl:  mediaUrl,
					EnclosurePath: filePath,
					EnclosureSize: podcastData.PodSizes[mediaIndex],
				}
				misc.ChannelLog("\t\t\t\t\t\t Go_deriveFilenames SENDING")
				mediaStream <- newMedia
				misc.ChannelLog("\t\t\t\t\t\t Go_deriveFilenames RECEIVED")
				limitFlag--
				if limitFlag == 0 {
					break limitCancel
				}
			} else {
				return 0, err
			}
		} else {
			haveCount++
			haveStr := strconv.Itoa(haveCount)
			globals.Console.Note("\t\t\t\tHave #" + haveStr + " " + nameOfFile + "\n")
		}
	}
	misc.ChannelLog("\t\t\t\t Go_deriveFilenames END") // never got here
	return possibleFiles, nil
}
