package testByNameOrUrl

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	test_helpers.DirRemove(progPath + "/By Name Or Url/")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
	rssData := map[string]string{
		"http://rss.ByNameOrUrl/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.ByNameOrUrl/file-1.ByNameOrUrl" length="42" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.ByNameOrUrl/file-2.ByNameOrUrl" length="43" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.ByNameOrUrl/file-1.ByNameOrUrl": `file 1 ByNameOrUrl `,
		"http://rss.ByNameOrUrl/file-2.ByNameOrUrl": `file 2 ByNameOrUrl`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.ByNameOrUrl", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `

 Adding 'By Name Or Url'
        
        
        Downloading 'By Name Or Url' podcast, 2 files, hit 's' to stop
        	file-2.ByNameOrUrl(read #0 43B)
        	file-1.ByNameOrUrl(read #0 42B)
        		 file-1.ByNameOrUrl (save #0, 0s)
						 		Size disparity, expected 42 bytes, but was 19
        		 file-2.ByNameOrUrl (save #0, 0s)
						 		Size disparity, expected 43 bytes, but was 18
`
const expectedAdds = `
Added 2 new files in 0s 
From http://rss.ByNameOrUrl/podcast.xml 
Into 'By Name Or Url'  
`

const expectedBads = ""

func Test_4_ByNameOrUrl(t *testing.T) { //AddByUrlAndName

	progBounds := setUp()

	cleanArgs := []string{"file-name.go", "http://rss.ByNameOrUrl/podcast.xml", "By", "Name", "Or", "Url"}
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, err := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, httpTest)
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

	podDir := progBounds.ProgPath + "/By Name Or Url"
	if _, err = os.Stat(podDir); err != nil {
		t.Fatal(podDir + " directory does not exist")
	}

	file1 := progBounds.ProgPath + "/By Name Or Url/file-1.ByNameOrUrl"
	if _, err = os.Stat(file1); err != nil {
		t.Fatal(file1 + " does not exist")
	}

}
