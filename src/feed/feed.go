package feed

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
)

func IsUrl(rssUrl string) bool {
	if strings.HasPrefix(rssUrl, "http") {
		return true
	}
	if strings.Contains(rssUrl, ".") && strings.Contains(rssUrl, "/") {
		return true
	}
	return false
}

func addHttp(rssUrl string) string {
	if strings.HasPrefix(rssUrl, "http") {
		return rssUrl
	}
	return "http://" + rssUrl
}

func IncGlobalCounters(incCounter *int) string {
	var mu sync.Mutex
	mu.Lock()
	*incCounter++
	countTemp := *incCounter
	mu.Unlock()
	countStr := fmt.Sprint(countTemp)
	return countStr
}

func ShowError(mediaUrl string) string {
	savedOut := fmt.Sprint("ERROR " + rss.NameOfFile(mediaUrl) + "\n")
	return savedOut
}

func ShowSizeError(expectedSize, writtenSize int) string {
	savedOut := ""
	if expectedSize > 0 {
		exSize := strconv.Itoa(expectedSize)
		wrSize := strconv.Itoa(writtenSize)
		savedOut = fmt.Sprint("\t\t\tSize disparity, expected " + exSize + " bytes, but was " + wrSize + "\n")
	}
	return savedOut
}

func ShowSaved(savedFiles *int, startProcess time.Time, mediaPath string) string {
	var roundTime time.Duration
	savedTemp := IncGlobalCounters(savedFiles)
	if !consts.IsTesting(os.Args) {
		sinceStart := time.Since(startProcess)
		roundTime = sinceStart.Round(time.Second) // NB if testing all times are 0s
	} else {
		savedTemp = "0"
	}
	saveNumMess := "(save #" + savedTemp + ", " + fmt.Sprint(roundTime) + ")"
	savedOut := fmt.Sprintf("\t\t %s %s", rss.NameOfFile(mediaPath), saveNumMess) + "\n"
	return savedOut
}

func PodcastName(progArgs []string) string {
	titleWords := make([]string, 0)
	for argInd := 1; argInd < len(progArgs); argInd++ {
		anArg := progArgs[argInd]
		if !IsUrl(anArg) {
			titleWords = append(titleWords, anArg)
		}
	}
	podcastTitle := strings.Join(titleWords, " ")
	return podcastTitle
}

func ShowProgress(fileEnc models.MediaEnclosure, readFiles *int) string {
	var fileCount string = "0"
	if !consts.IsTesting(os.Args) {
		fileCount = IncGlobalCounters(readFiles) // NB if testing all times are (read #0
	}
	mbGbLen := misc.GbOrMb(fileEnc.EnclosureSize)
	var readNumMess string
	if mbGbLen == "" {
		readNumMess = "(read #" + fileCount + ")\n"
	} else {
		readNumMess = "(read #" + fileCount + " " + mbGbLen + ")\n"
	}
	curMedia := "\t" + rss.NameOfFile(fileEnc.EnclosurePath) + readNumMess
	return curMedia
}

func ReadRss(rssUrl string, httpMedia models.HttpFn) ([]byte, error) {
	timeOut := misc.FileTimeout(globals.RssMaxReadFileTime)
	ctxRss, cancelRss := context.WithTimeout(context.Background(), timeOut)
	defer cancelRss()
	httpUrl := addHttp(rssUrl)
	rssResponse, err := httpMedia(ctxRss, httpUrl)
	if err != nil {
		return nil, err
	}
	if rssResponse.StatusCode != consts.HTTP_OK_RESP {
		return badStatusCode(httpUrl, rssResponse.Status, flaws.FLAW_E_80, err)
	}
	rssText, err := io.ReadAll(rssResponse.Body)
	if err != nil {
		return badReadAll(httpUrl, flaws.FLAW_E_81, err)
	}
	if len(rssText) == 0 {
		return readZero(httpUrl, flaws.FLAW_E_82, err)
	}
	return rssText, nil
}
