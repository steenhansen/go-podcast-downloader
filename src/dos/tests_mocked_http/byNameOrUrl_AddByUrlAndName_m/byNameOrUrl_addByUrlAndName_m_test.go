package testFinalMediaName

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"podcast-downloader/src/dos/console"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	test_helpers.DirRemove(progPath + "/by-name-or-url/")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string, numRetries int) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.FinalMediaName/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>

								<title>title tag</title>
								<itunes:title>changed tag</itunes:title>
								<item>
									<enclosure url="http://rss.FinalMediaName/not-this-name.txt" length="42" type="text/plain" />
								</item>
							
								<item>
									<enclosure url="http://rss.FinalMediaName/file-2.txt" length="97" type="text/plain" />
								</item>

								<item>
								<title>My Title2</title>
									<enclosure url="http://rss.FinalMediaName/file-Re-Direct.txt" length="63" type="text/plain" />
								</item>
							
								</channel>
						</rss>`,
		"http://rss.FinalMediaName/not-this-name.txt":  `file 1 FinalMediaName `,
		"http://rss.FinalMediaName/file-2.txt":         `file 2 FinalMediaName 01234567890`,
		"http://rss.FinalMediaName/file-Re-Direct.txt": `file 3 FinalRedirect 0123456789001234567890`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.FinalMediaName", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown FinalMediaName : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `
 Adding 'by-name-or-url'
        Downloading 'by-name-or-url' podcast, 3 files, hit 's' to stop
        	file-2.txt(read #0 97B)
        	not-this-name.txt(read #0 42B)
					file-Re-Direct.txt(read #0 63B)
        		 file-2.txt (save #0, 0s)
						 		Size disparity, expected 97 bytes, but was 33
        		 not-this-name.txt (save #0, 0s)
						 		Size disparity, expected 42 bytes, but was 22
						 file-Re-Direct.txt (save #0, 0s) 
						 		Size disparity, expected 63 bytes, but was 43
`
const expectedAdds = `
Added 3 new files in 0s 
From http://rss.FinalMediaName/podcast.xml 
Into 'by-name-or-url'  
`
const expectedBads = ""

func Test_2_ByNameOrUrl_AddByUrlAndName(t *testing.T) {
	progBounds := setUp()
	cleanArgs := []string{"file-name.go", "http://rss.FinalMediaName/podcast.xml", "by-name-or-url"}
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualAdds, podcastResults := console.ByNameOrUrl(cleanArgs, progBounds, keyStreamTest, httpTest)
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

	podDir := progBounds.ProgPath + "/by-name-or-url"
	if _, err = os.Stat(podDir); err != nil {
		t.Fatal(podDir + " directory does not exist")
	}

	file1 := progBounds.ProgPath + "/by-name-or-url/file-2.txt"
	if _, err = os.Stat(file1); err != nil {
		t.Fatal(err)
	}

	file2 := progBounds.ProgPath + "/by-name-or-url/file-Re-Direct.txt"
	if _, err = os.Stat(file2); err != nil {
		t.Fatal(err)
	}

	file3 := progBounds.ProgPath + "/by-name-or-url/not-this-name.txt"
	if _, err = os.Stat(file3); err != nil {
		t.Fatal(err)
	}

}
