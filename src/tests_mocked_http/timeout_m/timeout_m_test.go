package testTimeout

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	progBounds := test_helpers.TestBounds(progPath)
	progBounds.LoadOption = consts.LOW_LOAD
	globals.MediaMaxReadFileTime = time.Microsecond * 1
	return progBounds
}

func httpMocked(ctx context.Context, mediaUrl string) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.timeout-m/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.timeout-m/file-1.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-2.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-3.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-4.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-5.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-6.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-7.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-8.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-9.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.timeout-m/file-10.txt" length="16" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.timeout-m/file-1.txt":  `file 1 low-disk`,
		"http://rss.timeout-m/file-2.txt":  `file 2 low-disk`,
		"http://rss.timeout-m/file-3.txt":  `file 3 low-disk`,
		"http://rss.timeout-m/file-4.txt":  `file 4 low-disk`,
		"http://rss.timeout-m/file-5.txt":  `file 5 low-disk`,
		"http://rss.timeout-m/file-6.txt":  `file 6 low-disk`,
		"http://rss.timeout-m/file-7.txt":  `file 7 low-disk`,
		"http://rss.timeout-m/file-8.txt":  `file 8 low-disk`,
		"http://rss.timeout-m/file-9.txt":  `file 9 low-disk`,
		"http://rss.timeout-m/file-10.txt": `file 10 low-disk`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.timeout-m", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `
1 |  10 files |    0MB | timeout-m
         'Q' or a number + enter: Downloading 'timeout-m' podcast, 10 files, hit 's' to stop
        				Have #1 file-1.txt
        				Have #2 file-2.txt
        				Have #3 file-3.txt
        				Have #4 file-4.txt
        				Have #5 file-5.txt
        				Have #6 file-6.txt
        				Have #7 file-7.txt
        				Have #8 file-8.txt
        				Have #9 file-9.txt
        				Have #10 file-10.txt
`
const expectedAdds = `
No changes detected
`

func TestTimeout_m(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := menu.DisplayMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoiceNum("1"), httpMocked)
	var flawError flaws.FlawError
	err := podcastResults.SeriousError
	if errors.As(err, &flawError) {
		timeoutErr := flawError.Error()
		if timeoutErr != "Internet connection timed out by exceeding duration: 1µs" {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	actualConsole := globals.Console.All()
	expectedDiff := test_helpers.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualConsole), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedConsole))
	}

	if test_helpers.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(test_helpers.ClampActual(actualAdds), test_helpers.ClampExpected(expectedAdds))
	}

}
