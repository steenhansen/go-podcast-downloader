package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
)

// no title is ok, if user gives us a title

type xmlRssTitle struct {
	Title string `xml:"channel>title"`
}

type xmlItemTitles struct {
	Titles []string `xml:"channel>item>title"` // itunes:title are now ignored by changing them to itunes:tiXYZ
}

type xmlUrlLen struct {
	UrlKey string `xml:"url,attr"`
	LenKey string `xml:"length,attr"`
}

type xmlEnclosures struct {
	Enclosures []xmlUrlLen `xml:"channel>item>enclosure"`
}

func RssTitle(rssXml []byte) (string, error) {

	theChannel := xmlRssTitle{Title: ""}
	err := xml.Unmarshal([]byte(rssXml), &theChannel)
	if err != nil {
		return "", flaws.MissingTitle
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
		return 0, err
	}
	defer respMedia.Body.Close()
	if respMedia.StatusCode != 200 {
		return 0, flaws.BadContent.StartError(mediaUrl)
	}
	contentStr := ""
	mediaContent := make([]byte, 0)
	if !globals.EmptyFiles {
		mediaContent, err = io.ReadAll(respMedia.Body)
		if err != nil {
			return 0, flaws.BadContent.ContinueError(mediaUrl, err)
		}
		contentStr = string(mediaContent)
	}

	mediaFile, err := os.Create(filePath)
	if err != nil {
		return 0, flaws.CantCreateFileSerious.ContinueError(filePath, err)
	}
	defer mediaFile.Close()

	if strings.HasPrefix(contentStr, consts.HTML_404_BEGIN) {
		return 0, flaws.BadContent.StartError(mediaUrl)
	}
	err = misc.DiskPanic(len(mediaContent), minDiskMbs)
	if err != nil {
		return 0, err
	}
	writenBytes, err := mediaFile.Write(mediaContent)
	if err != nil {
		return 0, flaws.CantWriteFileSerious.ContinueError(filePath, err)
	}
	return writenBytes, nil
}

func HttpMedia(ctx context.Context, mediaUrl string) (*http.Response, error) {
	newReq, err := http.NewRequest(http.MethodGet, mediaUrl, nil)
	if err != nil {
		return nil, flaws.BadUrl.ContinueError(mediaUrl, err)
	}
	reqCtx := newReq.WithContext(ctx)
	httpClient := &http.Client{}
	respMedia, err := httpClient.Do(reqCtx)
	if err != nil {
		return nil, flaws.BadUrl.ContinueError(mediaUrl, err)
	}
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	return respMedia, nil
}

func FinalMediaName(ctx context.Context, mediaUrl string, httpMedia models.HttpFn) (string, error) {
	respMedia, err := httpMedia(ctx, mediaUrl)
	if err != nil {
		return "", nil
	}
	// https://stackoverflow.com/questions/16784419/in-golang-how-to-determine-the-final-url-after-a-series-of-redirects
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
