package feed

// feed

import (
	"context"
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
)

// siberiantimes.com/ecology/rss/
func IsUrl(url string) bool {
	if strings.HasPrefix(url, "http") {
		return true
	}
	if strings.Contains(url, ".") && strings.Contains(url, "/") {
		return true
	}
	return false
}

func rssUrl(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}
	return "http://" + url
}

// go run pod-down.go https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
func ReadRss(url string) ([]byte, error) { // 2 or 3
	url = rssUrl(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, flaws.BadUrl.ContinueError(url, err)
	}
	defer response.Body.Close()
	rss, err := io.ReadAll(response.Body)
	if err != nil || len(rss) == 0 {
		return nil, flaws.BadUrl.StartError(url)
	}
	return rss, nil
}

func DownloadAndWriteFile(ctx context.Context, furl, fpath string, min_disk_mbs int) error {
	req, err := http.NewRequest(http.MethodGet, furl, nil)
	if err != nil {
		return flaws.BadUrl.ContinueError(furl, err)
	}
	req = req.WithContext(ctx)
	c := &http.Client{}
	response, err := c.Do(req)

	if err != nil {
		if ctx.Err() == context.Canceled {
			return nil
		}
		return flaws.BadUrl.ContinueError(furl, err)
	}
	////////////////////////////////////
	if response.StatusCode != 200 {
		return flaws.BadContent.StartError(furl)
	}

	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		return nil
	}
	f, err := os.Create(fpath)
	if err != nil {
		return flaws.CantCreateFile.ContinueError(fpath, err)
	}
	defer f.Close()

	//////////////////// q*bert
	contentStr := string(content)
	if strings.HasPrefix(contentStr, consts.HTML_404_BEGIN) {
		return flaws.BadContent.StartError(furl)
	}

	err = misc.DiskPanic(len(content), min_disk_mbs)
	if err != nil {
		return err
	}
	_, err = f.Write(content)
	if err != nil {
		return flaws.CantWriteFile.ContinueError(fpath, err)
	}
	return nil
}

/* test cancel?

ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Millisecond*8000))
mediaStream := make(chan consts.UrlPathLength)
doneXStream := make(chan bool)

cancel()
downloaded, viewed := GoDownloadAndSaveFiles(ctx, mediaStream, doneXStream)
//should be 0,0

*/

func IncGlobalCounters(pont *int) string {
	var mu sync.Mutex
	mu.Lock()
	*pont++
	readTemp := *pont
	mu.Unlock()
	readStr := fmt.Sprint(readTemp)
	return readStr
}

func ShowError(mediaUrl string) string {
	savedOut := fmt.Sprint("ERROR " + misc.NameOfFile(mediaUrl))
	return savedOut
}

func ShowSaved(savedFiles *int, start time.Time, url string) string {
	var savedTemp string = "0"
	sinceStart := time.Since(start)
	var roundTime time.Duration
	if !misc.IsTesting(os.Args) {
		savedTemp = IncGlobalCounters(savedFiles)
		roundTime = sinceStart.Round(time.Second) // NB if testing all times are 0s
	}
	saveNumMess := "(save #" + savedTemp + ", " + fmt.Sprint(roundTime) + ")"
	savedOut := fmt.Sprintf("\t\t %s %s", misc.NameOfFile(url), saveNumMess)
	return savedOut
}

// ./pd-console NASA, image, of, the, day => "NASA image of the day"
func PodcastName(Args []string) string {
	names := make([]string, 0)
	for i := 1; i < len(Args); i++ {
		argument := Args[i]
		if !IsUrl(argument) {
			names = append(names, argument)
		}
	}
	name := strings.Join(names, " ")
	return name
}

func ShowProgress(media consts.UrlPathLength, readFiles *int) string {
	var readTemp string = "0"
	if !misc.IsTesting(os.Args) {
		readTemp = IncGlobalCounters(readFiles) // NB if testing all times are (read #0
	}

	mbGbLen := misc.GbOrMb(media.Length)
	var readNumMess string
	if mbGbLen == "" {
		readNumMess = "(read #" + readTemp + ")"
	} else {
		readNumMess = "(read #" + readTemp + " " + mbGbLen + ")"
	}
	progress := "\t" + misc.NameOfFile(media.Url) + readNumMess
	return progress
}
