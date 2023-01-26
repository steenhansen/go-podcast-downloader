package media

import (
	"context"
	"os"
	"regexp"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
)

func CurDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

func dirTitle(podTitle, rssUrl string) string {
	if podTitle == "" {
		podTitle = rssUrl
	}
	xmlEscaped := regexp.MustCompile(`&[^;]*;`) // &amp
	safeXml := xmlEscaped.ReplaceAllLiteralString(podTitle, "")

	multSpaces := regexp.MustCompile(`\s\s+`)
	singleSpaces := multSpaces.ReplaceAllLiteralString(safeXml, " ")

	alphaRe := regexp.MustCompile(`[^a-zA-Z0-9 _-]+`)
	safeTitle := alphaRe.ReplaceAllLiteralString(singleSpaces, "")
	return safeTitle
}

func InitFolder(progPath, podTitle, rssUrl string) (string, error) {
	safeTitle := dirTitle(podTitle, rssUrl)
	containDir := progPath + "/" + safeTitle
	if _, err := os.Stat(containDir); os.IsNotExist(err) {
		if err := os.Mkdir(containDir, os.ModePerm); err != nil {
			return "", flaws.CantCreateDir.ContinueError(containDir, err)
		}
	}
	originRss := containDir + "/" + consts.URL_OF_RSS
	f, err := os.Create(originRss)
	if err != nil {
		return "", flaws.CantCreateFile.ContinueError(originRss, err)
	}
	defer f.Close()
	_, err = f.Write([]byte(rssUrl))
	if err != nil {
		return "", flaws.CantWriteFile.ContinueError(originRss, err)
	}
	return containDir, nil
}

func SaveDownloadedMedia(ctx context.Context, podcastData consts.PodcastData, mediaStream chan<- consts.UrlPathLength, limitFlag int) (int, string, error) {
	varieties := misc.VarietiesSet{}
	possibleFiles := 0
limitCancel:
	for i, furl := range podcastData.Medias {
		select {
		case <-ctx.Done():
			break limitCancel
		default:
			possibleFiles++
			fname := misc.NameOfFile(furl)
			fpath := podcastData.MediaPath + "/" + fname
			_, err := os.Stat(fpath)
			if err != nil {
				if os.IsNotExist(err) {
					newMedia := consts.UrlPathLength{Url: furl,
						Path:   fpath,
						Length: podcastData.Lengths[i],
					}
					mediaStream <- newMedia
					varieties.AddVariety(fname)
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
	fTypes := varieties.VarietiesString(" ")
	return possibleFiles, fTypes, nil
}
