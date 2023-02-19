package rss

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
)

//  https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/bad-url/does-not-exist.rss

const URL2 = consts.TEST_DIR_URL + "bad-url/does-not-exist.rss"

//var expected_err = flaws.InvalidRssURL.MakeFlaw(URL2)

func httpMedia(ctx context.Context, mediaUrl string) (*http.Response, error) {
	return nil, flaws.InvalidRssURL.MakeFlaw(URL2)
}

func TestInvalidXml(t *testing.T) {
	keyStream := make(chan string)
	_, _, _, _, err := podcasts.ReadRssUrl(URL2, httpMedia, keyStream)
	if !errors.Is(err, flaws.InvalidRssURL) {
		t.Fatal(err)
	}
}
