package tests_mocked_http

import (
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
)

// go run pod-down.go http://127.0.0.1:9000/boom
// Get "http://127.0.0.1:9000/boom": dial tcp 127.0.0.1:9000: connectex: No connection could be made because the target machine actively refused it.

// go run pod-down.go https://www.naxxxsa.gov/rss/dyn/lg_image_of_the_dayXX.rss
// Get "https://www.naxxxsa.gov/rss/dyn/lg_image_of_the_dayXX.rss": dial tcp: lookup www.naxxxsa.gov: no such host

// go run pod-down.go https://raw.githubusercontent.com/steenhansen/projects/main/images/zero-bytes.txt
// emtpy rss file https://raw.githubusercontent.com/steenhansen/projects/main/images/zero-bytes.txt

func Test_6_IsUrl(t *testing.T) {
	okUrl := "siberiantimes.com/ecology/rss/"
	isUrl := feed.IsUrl(okUrl)
	if !isUrl {
		t.Fatal(`TestIsUrl A failed`)
	}
	badUrl1 := "siberiantimes.com"
	isUrl = feed.IsUrl(badUrl1)
	if isUrl {
		t.Fatal(`TestIsUrl B failed`)
	}
	badUrl2 := "siberiantimescom/"
	isUrl = feed.IsUrl(badUrl2)
	if isUrl {
		t.Fatal(`TestIsUrl C failed`)
	}
}

// func httpMedia(ctx context.Context, mediaUrl string) (*http.Response, error) {
// 	body := `<?xml version="1.0" encoding="UTF-8"?>
// <rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
//   <channel>
//     <title>title tag</title>
//     <item>
//       <enclosure url="https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/not-missing.txt" length="11" type="text/plain" />
//     </item>
//     <item>
//       <enclosure url="https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/no-such-file.txt" length="12" type="text/plain" />
//     </item>
//   </channel>
// </rss>`
// 	t := &http.Response{
// 		Status:        "200 OK",
// 		StatusCode:    200,
// 		Proto:         "HTTP/1.1",
// 		ProtoMajor:    1,
// 		ProtoMinor:    1,
// 		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
// 		ContentLength: int64(len(body)),
// 		Header:        make(http.Header, 0),
// 	}
// 	return t, nil
// }

// func TestReadRss(t *testing.T) {
// 	badUrl := "doesnot.exist/"
// 	rssXml, rssFiles, rssSizes, err := podcasts.ReadRssUrl(badUrl, httpMedia)
// 	fmt.Println("rssXml", rssXml)
// 	fmt.Println("rssFiles", rssFiles)
// 	fmt.Println("rssSizes", rssSizes)
// 	if !errors.Is(err, flaws.BadUrl) {
// 		t.Fatal(`TestReadRss A failed`)
// 	}
// }
