package testDisplayMenu

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/menu"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

func setUp() consts.ProgBounds {
	progPath := misc.CurDir()
	os.Remove(progPath + "/Display Menu/file-5.txt")
	os.Remove(progPath + "/Display Menu/file-6.txt")
	progBounds := testings.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
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
		httpResp := testings.Http200Resp("rss.DisplayMenu", thePath, theData)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `

 1 |              txt |   2 files |    0MB | Display Menu
 'Q' or a number + enter: Downloading 'Display Menu' podcast, 4 files, hit 's' to stop
	file-5.txt(read #0 44B)
	file-6.txt(read #0 45B)
		 file-5.txt (save #0, 0s)
			Size disparity, expected 44 bytes, but was 19
		 file-6.txt (save #0, 0s)
			Size disparity, expected 45 bytes, but was 18
`
const expectedAdds = `
Added 2 new 'txt' file(s) in 0s 
From http://rss.DisplayMenu/podcast.xml 
Into 'Display Menu'
`

const expectedBads = ""

func Test_DisplayMenu(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, err := menu.DisplayMenu(progBounds, keyStream, testings.KeyboardMenuChoiceNum("1"), httpTest)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()

	expectedDiff := testings.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(testings.ClampActual(actualConsole), testings.ClampMapDiff(expectedDiff), testings.ClampExpected(expectedConsole))
	}

	if testings.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(testings.ClampActual(actualAdds), testings.ClampExpected(expectedAdds))
	}

	if testings.NotSameTrimmed(actualBads, expectedBads) {
		t.Fatal(testings.ClampActual(actualBads), testings.ClampExpected(expectedBads))
	}

	file5 := progBounds.ProgPath + "/Display Menu/file-5.txt"
	if _, err = os.Stat(file5); err != nil {
		t.Fatal(file5 + " does not exist")
	}
}
