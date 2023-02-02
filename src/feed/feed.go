package feed

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
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

func ReadRss(rssUrl string) ([]byte, error) {
	httpUrl := addHttp(rssUrl)
	rssResponse, err := http.Get(httpUrl)
	if err != nil {
		return nil, flaws.BadUrl.ContinueError(httpUrl, err)
	}
	defer rssResponse.Body.Close()
	rssText, err := io.ReadAll(rssResponse.Body)
	if err != nil {
		return nil, flaws.BadUrl.ContinueError(httpUrl, err)
	}
	if len(rssText) == 0 {
		return nil, flaws.EmptyRss.StartError(httpUrl)
	}
	return rssText, nil
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
	savedOut := fmt.Sprint("ERROR " + rss.NameOfFile(mediaUrl))
	return savedOut
}

func ShowSaved(savedFiles *int, startProcess time.Time, mediaPath string) string {
	//var savedTemp string = "0"
	var roundTime time.Duration
	savedTemp := IncGlobalCounters(savedFiles)
	if !misc.IsTesting(os.Args) {
		//savedTemp = IncGlobalCounters(savedFiles)
		sinceStart := time.Since(startProcess)
		roundTime = sinceStart.Round(time.Second) // NB if testing all times are 0s
	} else {
		savedTemp = "0"
	}
	saveNumMess := "(save #" + savedTemp + ", " + fmt.Sprint(roundTime) + ")"
	savedOut := fmt.Sprintf("\t\t %s %s", rss.NameOfFile(mediaPath), saveNumMess)
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

func ShowProgress(fileEnc consts.MediaEnclosure, readFiles *int) string {
	var fileCount string = "0"
	if !misc.IsTesting(os.Args) {
		fileCount = IncGlobalCounters(readFiles) // NB if testing all times are (read #0
	}
	mbGbLen := misc.GbOrMb(fileEnc.EnclosureSize)
	var readNumMess string
	if mbGbLen == "" {
		readNumMess = "(read #" + fileCount + ")"
	} else {
		readNumMess = "(read #" + fileCount + " " + mbGbLen + ")"
	}
	curMedia := "\t" + rss.NameOfFile(fileEnc.EnclosurePath) + readNumMess
	return curMedia
}
