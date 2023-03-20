package media

import (
	"context"
	"os"
	"regexp"
	"strconv"
	"strings"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
)

func DirTitle(podTitle, rssUrl string) string {
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
	forcingTitle := globals.ForceTitle
	return ReSaveFolder(forcingTitle, progPath, podTitle, rssUrl)
}

func ReSaveFolder(forcingTitle bool, progPath, podTitle, rssUrl string) (string, bool, error) {
	safeTitle := DirTitle(podTitle, rssUrl)
	containDir := progPath + "/" + safeTitle
	dirNotExist := false
	if _, err := os.Stat(containDir); os.IsNotExist(err) {
		if err := os.Mkdir(containDir, os.ModePerm); err != nil {
			return DoesntExist(containDir, flaws.FLAW_E_70, err)
		}
		dirNotExist = true
	}
	originRss := containDir + "/" + consts.URL_OF_RSS_FN
	rssAddrFile, err := os.Create(originRss)
	if err != nil {
		return CannotCreate(originRss, flaws.FLAW_E_71, err)
	}
	defer rssAddrFile.Close()
	_, err = rssAddrFile.Write([]byte(rssUrl)) // media files change size on ad injection
	if err != nil {
		return WriteError(originRss, flaws.FLAW_E_72, err)
	}
	if forcingTitle {
		rssAddrFile.Write([]byte("\n" + consts.OPTION_FORCE_TITLE))
	}
	return containDir, dirNotExist, nil
}

func TitleToName(podTitle string) string {
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
	return trimmedName
}

func FileExten(finalFileName string) string {
	filePieces := strings.Split(finalFileName, ".")
	fileExt := filePieces[len(filePieces)-1]
	return fileExt
}

func chooseName(finalFileName string, podPath string, podTitle string, podcastData models.PodcastData) (filePath string) {
	filePath = podPath + "/"
	fileExt := FileExten(finalFileName)
	trimmedName := TitleToName(podTitle)
	////////////
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
