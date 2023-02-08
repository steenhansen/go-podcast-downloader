package rss

import (
	"context"
	"net/http"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
)

//  https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/internet-tests/bad-url/does-not-exist.rss

const URL2 = consts.TEST_DIR_URL + "bad-url/does-not-exist.rss"

var expected_err = flaws.BadUrl.StartError(URL2)

func httpMedia(ctx context.Context, mediaUrl string) (*http.Response, error) {
	return nil, flaws.BadUrl.StartError(URL2)
}

func TestInvalidXml(t *testing.T) {
	_, _, _, _, err := podcasts.ReadRssUrl(URL2, httpMedia)
	if err != expected_err {
		t.Fatal(err)
	}
}
