package rss

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
)

// no title is ok, if user gives us a title
func RssTitle(rssXml []byte) (string, error) {
	type XmlTitle struct {
		Title string `xml:"channel>title"`
	}
	theChannel := XmlTitle{Title: ""}
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

func RssItems(rss []byte) ([]string, []int, error) {
	type XmlAttrib struct {
		UrlKey string `xml:"url,attr"`
		LenKey string `xml:"length,attr"`
	}
	type XmlEnclosures struct {
		Enclosures []XmlAttrib `xml:"channel>item>enclosure"`
	}
	enclosures := XmlEnclosures{}
	err := xml.Unmarshal([]byte(rss), &enclosures)
	if err != nil {
		return nil, nil, err
	}
	mediaUrls := make([]string, len(enclosures.Enclosures))
	mediaSizes := make([]int, len(enclosures.Enclosures))
	if len(mediaUrls) == 0 {
		return nil, nil, flaws.EmptyItems
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
	return mediaUrls, mediaSizes, nil
}

func NameOfFile(mediaUrl string) string {
	urlParts := strings.Split(mediaUrl, "/")
	fileNameQuery := urlParts[len(urlParts)-1]
	fileNameNoQuery := strings.Split(fileNameQuery, "?")
	return fileNameNoQuery[0]
}

func DownloadAndWriteFile(ctx context.Context, mediaUrl, filePath string, minDiskMbs int, httpMedia consts.HttpFunc) (int, error) {
	respMedia, err := httpMedia(ctx, mediaUrl)
	if err != nil {
		return 0, err
	}
	defer respMedia.Body.Close()
	if respMedia.StatusCode != 200 {
		return 0, flaws.BadContent.StartError(mediaUrl)
	}
	mediaContent, err := io.ReadAll(respMedia.Body)
	if err != nil {
		return 0, flaws.BadContent.ContinueError(mediaUrl, err)
	}
	mediaFile, err := os.Create(filePath)
	if err != nil {
		return 0, flaws.CantCreateFileSerious.ContinueError(filePath, err)
	}
	defer mediaFile.Close()
	contentStr := string(mediaContent)
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

func FinalMediaName(ctx context.Context, mediaUrl string, httpMedia consts.HttpFunc) (string, error) {
	respMedia, err := httpMedia(ctx, mediaUrl)
	if err != nil {
		return "", nil
	}
	if respMedia.StatusCode != consts.HTTP_OK_RESP {
		missingFinalName := NameOfFile(mediaUrl)
		return missingFinalName, nil
	}
	finalQueried := respMedia.Request.URL.String()
	finalFileName := NameOfFile(finalQueried)
	return finalFileName, nil
}
