package testDisplayMenu

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/console"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	os.Remove(progPath + "/display-menu-m/file-5.txt")
	os.Remove(progPath + "/display-menu-m/file-6.txt")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string, numRetries int) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.DisplayMenu/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.DisplayMenu/file-3.txt" length="42" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.DisplayMenu/file-4.txt" length="43" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.DisplayMenu/file-5.txt" length="44" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.DisplayMenu/file-6.txt" length="45" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.DisplayMenu/file-3.txt": `file 3 DisplayMenu `,
		"http://rss.DisplayMenu/file-4.txt": `file 4 DisplayMenu`,
		"http://rss.DisplayMenu/file-5.txt": `file 5 DisplayMenu `,
		"http://rss.DisplayMenu/file-6.txt": `file 6 DisplayMenu`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.DisplayMenu", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `

1 |   2 files |    0MB | display-menu-m
 'Q' or a number + enter: Downloading 'display-menu-m' podcast, 4 files, hit 's' to stop
	file-5.txt(read #0 44B)
	file-6.txt(read #0 45B)
		 file-5.txt (save #0, 0s)
		 		Size disparity, expected 44 bytes, but was 19
		 file-6.txt (save #0, 0s) 
		 		Size disparity, expected 45 bytes, but was 18
`
const expectedAdds = `
Added 2 new files in 0s 
From http://rss.DisplayMenu/podcast.xml 
Into 'display-menu-m'
`

const expectedBads = ""

func TestDisplayMenu_m(t *testing.T) {
	progBounds := setUp()
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := console.DisplayMenu(progBounds, keyStreamTest, test_helpers.KeyboardMenuChoiceNum("1"), httpTest)
	err := podcastResults.SeriousError
	if err != nil {
		t.Fatal(err)
	}
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()

	expectedDiff := test_helpers.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualConsole), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedConsole))
	}

	if test_helpers.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(test_helpers.ClampActual(actualAdds), test_helpers.ClampExpected(expectedAdds))
	}

	if test_helpers.NotSameTrimmed(actualBads, expectedBads) {
		t.Fatal(test_helpers.ClampActual(actualBads), test_helpers.ClampExpected(expectedBads))
	}

	file5 := progBounds.ProgPath + "/display-menu-m/file-5.txt"
	if _, err = os.Stat(file5); err != nil {
		t.Fatal(file5 + " does not exist")
	}
}
