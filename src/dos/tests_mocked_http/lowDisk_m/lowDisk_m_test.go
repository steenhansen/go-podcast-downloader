package testLowDisk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"podcast-downloader/src/dos/console"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	progBounds := test_helpers.TestBounds(progPath)
	progBounds.MinDisk = 1_000_000_000_000_000
	return progBounds
}

func httpMocked(ctx context.Context, mediaUrl string, numRetries int) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.Low-Disk/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.Low-Disk/file-1.txt" length="16" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Low-Disk/file-2.txt" length="16" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.Low-Disk/file-1.txt": `file 1 low-disk`,
		"http://rss.Low-Disk/file-2.txt": `file 2 low-disk`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.Low-Disk", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `
       1 |   1 files |    0MB | low-disk-m
         'Q' or a number + enter: Downloading 'low-disk-m' podcast, 2 files, hit 's' to stop
                                        Have #1 file-1.txt
                file-2.txt(read #0 16B)
        ERROR: E_15
        low disk space, 22GB free, need minimum 909TB to proceed FILE: file-2.txt
`
const expectedAdds = `
No changes detected
`

const expectedBads = `
 E_15 : low disk space, xxGB free, need minimum 909TB to proceed
`

func TestLowDisk_m(t *testing.T) {
	progBounds := setUp()
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := console.DisplayMenu(progBounds, keyStreamTest, test_helpers.KeyboardMenuChoiceNum("1"), httpMocked)
	var flawError flaws.FlawError
	err := podcastResults.SeriousError
	fmt.Println("errror XXXXXXXXXXXXXXX", err, "YYYYYYYYYYYYYYYYYYYYYYYYYYYYY")
	fmt.Println("errror ccccccccccccc", err.Error(), "ddddddddddddd")
	if errors.As(err, &flawError) {
		lowErr := flawError.Error()
		safeErr := test_helpers.ReplaceXxGbFree(lowErr)
		if safeErr != "low disk space, xxGB free, need minimum 909TB to proceed" {
			t.Fatal(err)
		}
	} else {
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

	actualSafe := test_helpers.ReplaceXxGbFree(actualBads)
	expectedSafe := test_helpers.ReplaceXxGbFree(expectedBads)
	if test_helpers.NotSameTrimmed(actualSafe, expectedSafe) {
		t.Fatal(test_helpers.ClampActual(actualSafe), test_helpers.ClampExpected(expectedSafe))
	}

}
