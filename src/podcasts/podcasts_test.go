package podcasts

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

var expectedFiles = []string{
	"http://rss.ReadRssUrl/not-missing.ReadRssUrl",
	"http://rss.ReadRssUrl/no-such-file.ReadRssUrl"}

var expectedSizes = []int{11, 12}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
	rssData := map[string]string{
		"http://rss.ReadRssUrl/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.ReadRssUrl/not-missing.ReadRssUrl" length="11" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.ReadRssUrl/no-such-file.ReadRssUrl" length="12" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.ReadRssUrl/not-missing.ReadRssUrl":  `not missing `,
		"http://rss.ReadRssUrl/no-such-file.ReadRssUrl": `no such file`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		httpResp := testings.Http200Resp("rss.ReadRssUrl", thePath, theData)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

func Test_ReadRssUrl(t *testing.T) {
	_, actualFiles, actualSizes, _ := ReadRssUrl("http://rss.ReadRssUrl/podcast.xml", httpTest)
	if !reflect.DeepEqual(actualFiles, expectedFiles) {
		actualJoin := "\n" + strings.Join(actualFiles, "\n")
		expectedJoin := "\n" + strings.Join(expectedFiles, "\n")
		t.Fatal(testings.ClampActual(actualJoin), testings.ClampExpected(expectedJoin))
	}
	if !reflect.DeepEqual(actualSizes, expectedSizes) {
		t.Fatal(actualSizes, expectedSizes)
	}
}
