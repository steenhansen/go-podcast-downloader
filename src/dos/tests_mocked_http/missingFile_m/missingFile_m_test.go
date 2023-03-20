package t1

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/terminal"
	"podcast-downloader/src/dos/test_helpers"
)

/*

https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_mock_internet/missingFile_m/git-server-source/missing-file-r.rss

*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	os.Remove(progPath + "/missing-file-m/not-missing.txt")
	progBounds := test_helpers.TestBounds(progPath)

	globals.MediaMaxReadFileTime = time.Second * 5

	return progBounds
}

func httpMocked(ctx context.Context, mediaUrl string, numRetries int) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.Missing-File/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.Missing-File/not-missing.txt" length="15" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Missing-File/no-such-file.txt" length="15" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.Missing-File/not-missing.txt": `file 1 low-disk`,
		//"http://rss.Missing-File/no-such-file.txt": `file 2 low-disk`,
	}
	if mediaUrl == "http://rss.Missing-File/no-such-file.txt" {
		_, err := rss.Not200Flaw("404", mediaUrl, flaws.FLAW_E_10)
		return nil, err
	}
	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.Missing-File", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedMenu string = `
1 |   0 files |    0MB | missing-file-m
 'Q' or a number + enter:
`

const expectedConsole string = `
Downloading 'missing-file-m' podcast, 2 files, hit 's' to stop
        	not-missing.txt(read #0 15B)
        		 not-missing.txt (save #0, 0s)
        			Size disparity, expected 15 bytes, but was 15
`

const expectedAdds = `
Added 1 new files in 0s
From http://rss.Missing-File/podcast.xml
Into 'missing-file-m'
`

const expectedBads = `	
E_10
HTTP error 404 : http://rss.Missing-File/no-such-file.txt
`

func TestMissingFile_m(t *testing.T) {
	progBounds := setUp()
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	globals.Console.Clear()
	getMenuChoice := test_helpers.KeyboardMenuChoice_1
	actualAdds, podcastResults := terminal.AfterMenu(progBounds, keyStreamTest, getMenuChoice, httpMocked)
	fmt.Println("my podcastResults", podcastResults.SeriousError.Error())
	fmt.Println("my podcastResults")

	actualConsole := globals.Console.All()

	actualBads := podcastResults.SeriousError.Error()

	if test_helpers.NotSameTrimmed(actualMenu, expectedMenu) {
		t.Fatal(test_helpers.ClampActual(actualMenu), test_helpers.ClampExpected(expectedMenu))
	}

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
}
