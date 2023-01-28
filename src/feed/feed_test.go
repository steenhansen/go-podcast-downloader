package feed

import (
	"errors"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

// go run pod-down.go http://127.0.0.1:9000/boom
// Get "http://127.0.0.1:9000/boom": dial tcp 127.0.0.1:9000: connectex: No connection could be made because the target machine actively refused it.

// go run pod-down.go https://www.naxxxsa.gov/rss/dyn/lg_image_of_the_dayXX.rss
// Get "https://www.naxxxsa.gov/rss/dyn/lg_image_of_the_dayXX.rss": dial tcp: lookup www.naxxxsa.gov: no such host

// go run pod-down.go https://raw.githubusercontent.com/steenhansen/projects/main/images/zero-bytes.txt
// emtpy rss file https://raw.githubusercontent.com/steenhansen/projects/main/images/zero-bytes.txt

func TestIsUrl(t *testing.T) {
	okUrl := "siberiantimes.com/ecology/rss/"
	isUrl := IsUrl(okUrl)
	if !isUrl {
		t.Fatalf(`TestIsUrl A failed`)
	}
	badUrl1 := "siberiantimes.com"
	isUrl = IsUrl(badUrl1)
	if isUrl {
		t.Fatalf(`TestIsUrl B failed`)
	}
	badUrl2 := "siberiantimescom/"
	isUrl = IsUrl(badUrl2)
	if isUrl {
		t.Fatalf(`TestIsUrl C failed`)
	}
}

func TestReadRss(t *testing.T) {
	badUrl := "doesnot.exist/"
	_, err := ReadRss(badUrl)
	if !errors.Is(err, flaws.BadUrl) {
		t.Fatalf(`TestReadRss A failed`)
	}
}
