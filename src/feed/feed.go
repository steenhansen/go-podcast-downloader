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

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
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

/* test cancel?

ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*8000))
mediaStream := make(chan consts.MediaEnclosure)
doneXStream := make(chan bool)

cancel()
downloaded, viewed := GoDownloadAndSaveFiles(ctx, mediaStream, doneXStream)
//should be 0,0

*/

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
	if !globals.EmptyFiles {
		exSize := strconv.Itoa(expectedSize)
		wrSize := strconv.Itoa(writtenSize)
		savedOut = fmt.Sprint("\t\t\tSize disparity, expected " + exSize + " bytes, but was " + wrSize + "\n")
	}
	return savedOut
}

func ShowSaved(savedFiles *int, startProcess time.Time, mediaPath string) string {
	var roundTime time.Duration
	savedTemp := IncGlobalCounters(savedFiles)
	if !misc.IsTesting(os.Args) {
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
	if !misc.IsTesting(os.Args) {
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
	ctxRss, cancelRss := context.WithTimeout(context.Background(), time.Duration(consts.MAX_READ_FILE_TIME))
	defer cancelRss()
	httpUrl := addHttp(rssUrl)
	rssResponse, err := httpMedia(ctxRss, httpUrl)
	if err != nil {
		return nil, err
	}
	defer rssResponse.Body.Close()
	if rssResponse.StatusCode != consts.HTTP_OK_RESP {
		return nil, flaws.BadUrl.StartError(httpUrl)
	}
	rssText, err := io.ReadAll(rssResponse.Body)
	if err != nil {

		return nil, flaws.BadUrl.ContinueError(httpUrl, err)
	}
	if len(rssText) == 0 {
		return nil, flaws.EmptyRss.StartError(httpUrl)
	}
	return rssText, nil
}
