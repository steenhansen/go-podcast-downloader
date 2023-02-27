package rss

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
)

func RssTitle(rssXml []byte) (string, error) {

	theChannel := xmlRssTitle{Title: ""}
	err := xml.Unmarshal([]byte(rssXml), &theChannel)
	if err != nil {
		return "", flaws.InvalidXmlTitle
	}
	title := strings.TrimSpace(theChannel.Title)
	if len(title) == 0 {
		return "", flaws.EmptyTitle
	}
	return theChannel.Title, nil
}

func RssItems(orgRss []byte) ([]string, []string, []int, error) {
	var re1 = regexp.MustCompile(`:title\>`)                          // <itunes:title>itunes title</itunes:title>
	noItunesTitles := re1.ReplaceAllString(string(orgRss), ":tiXYZ>") // <itunes:TEMP_TITLE>itunes title</itunes:TEMP_TITLE>

	theTitles := xmlItemTitles{}
	xml.Unmarshal([]byte(noItunesTitles), &theTitles)
	mediaTitles := make([]string, len(theTitles.Titles))
	for i, itemTitle := range theTitles.Titles {
		mediaTitles[i] = string(itemTitle)
	}

	enclosures := xmlEnclosures{}
	err := xml.Unmarshal([]byte(noItunesTitles), &enclosures)
	if err != nil {
		return nil, nil, nil, err
	}
	mediaUrls := make([]string, len(enclosures.Enclosures))
	mediaSizes := make([]int, len(enclosures.Enclosures))

	if len(mediaUrls) == 0 {
		return nil, nil, nil, flaws.EmptyItems
	}
	for i, v := range enclosures.Enclosures {
		mediaUrls[i] = string(v.UrlKey)
		size, err := strconv.Atoi(v.LenKey)
		if err != nil {
			mediaSizes[i] = 0
		} else {
			mediaSizes[i] = size
		}
	}
	return mediaTitles, mediaUrls, mediaSizes, nil
}

func NameOfFile(mediaUrl string) string {
	urlParts := strings.Split(mediaUrl, "/")
	fileNameQuery := urlParts[len(urlParts)-1]
	fileNameNoQuery := strings.Split(fileNameQuery, "?")
	return fileNameNoQuery[0]
}

func DownloadAndWriteFile(ctx context.Context, mediaUrl, filePath string, minDiskMbs int, httpMedia models.HttpFn) (int, error) {
	respMedia, err := httpMedia(ctx, mediaUrl)
	if err != nil {
		return 0, err // NB dns errors such as "no such host" from from 'Breaking Points' after retries
	}
	if respMedia.StatusCode != consts.HTTP_OK_RESP {
		return not200Flaw(respMedia.Status, mediaUrl, flaws.FLAW_E_10)
	}
	mediaContent := make([]byte, 0)
	if !globals.EmptyFilesTest {
		mediaContent, err = io.ReadAll(respMedia.Body)
		if err != nil {
			return readAllFlaw(mediaUrl, flaws.FLAW_E_11, err)
		}
	}
	contentStr := string(mediaContent)

	if strings.HasPrefix(contentStr, consts.HTML_404_BEGIN) {
		return was404Flaw(filePath, mediaUrl, flaws.FLAW_E_13, err)
	}

	mediaFile, err := os.Create(filePath)
	if err != nil {
		return osCreateFlaw(filePath, flaws.FLAW_E_12, err)
	}

	err = misc.DiskPanic(len(mediaContent), minDiskMbs)
	if err != nil {
		return diskPanicFlaw(mediaFile, filePath, flaws.FLAW_E_15, err)
	}
	writtenBytes, err := mediaFile.Write(mediaContent)
	if err != nil {
		return badWriteFlaw(mediaFile, filePath, flaws.FLAW_E_14, err)
	}
	if writtenBytes < 1 {
		return length0Flaw(mediaFile, filePath, flaws.FLAW_E_16)
	}
	if !globals.EmptyFilesTest && writtenBytes != len(mediaContent) {
		return lengthWrongFlaw(mediaFile, filePath, flaws.FLAW_E_17, writtenBytes, len(mediaContent))
	}
	mediaFile.Close()
	return writtenBytes, nil
}

func retryHttp(ctx context.Context, tryHttpMedia func() (*http.Response, error)) (respMedia *http.Response, err error) {
	sleepTime := consts.RETRY_SLEEP_START
	var retry int
	var dnsError *net.DNSError
	for retry = 0; retry < consts.HTTP_RETRIES; retry++ {
		if retry > 0 {
			time.Sleep(time.Duration(sleepTime) * time.Second)
			sleepTime *= 2
		}
		respMedia, err = tryHttpMedia()
		if retry == 0 && globals.DnsErrorsTest {
			noSuchHost := &net.DNSError{Err: "no such host TEST"}
			err = noSuchHost
		}
		if ctx.Err() == context.Canceled {
			return nil, context.Canceled
		}
		if err == nil {
			return respMedia, nil
		}
		dnsMess := ""

		if errors.As(err, &dnsError) {
			dnsMess = dnsError.Err
		}
		globals.Console.Note("Retrying " + strconv.Itoa(retry+1) + consts.ERROR_SEPARATOR + dnsMess + "\n")
	}
	retryCount := strconv.Itoa(retry)
	reqUrl := respMedia.Request.URL
	mediaUrl := reqUrl.Scheme + "://" + reqUrl.Host + reqUrl.Path
	retryErrMess := retryCount + ", for " + mediaUrl + " " + err.Error()
	return badRetryHttp(retryErrMess, flaws.FLAW_E_40, err)
}

func callHttpMedia(ctx context.Context, mediaUrl string) (*http.Response, error) {
	newReq, err := http.NewRequest(http.MethodGet, mediaUrl, nil)
	if err != nil {
		return badCallHttp(mediaUrl, flaws.FLAW_E_30, err)
	}
	reqCtx := newReq.WithContext(ctx)
	httpClient := &http.Client{}
	respMedia, err := httpClient.Do(reqCtx)
	return respMedia, err
}

func HttpReal(ctx context.Context, mediaUrl string) (*http.Response, error) {
	respMedia, err := retryHttp(ctx, func() (*http.Response, error) { return callHttpMedia(ctx, mediaUrl) })
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	if err != nil {
		return badHttp(mediaUrl, flaws.FLAW_E_20, err)
	}
	return respMedia, nil
}

/*
BBC News Top stories  go run ./ --emptyFiles podcasts.files.bbci.co.uk/p02nq0gn.rss
Witness History       go run ./ --emptyFiles podcasts.files.bbci.co.uk/p004t1hd.rss
*/
func FinalMediaName(ctx context.Context, mediaUrl string, httpMedia models.HttpFn) (string, error) {
	respMedia, err := httpMedia(ctx, mediaUrl)
	if err != nil {
		return "", err
	}
	finalQueried := respMedia.Request.URL.String()
	contentDisposition := respMedia.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		contentLines := strings.Split(contentDisposition, "\"")
		if len(contentLines) > 2 {
			contentFilename := contentLines[1] // podcasts.files.bbci.co.uk/p02nq0gn.rss
			return contentFilename, nil        // podcasts.files.bbci.co.uk/p004t1hd.rss
		}
	}
	finalFileName := NameOfFile(finalQueried)
	return finalFileName, nil
}
