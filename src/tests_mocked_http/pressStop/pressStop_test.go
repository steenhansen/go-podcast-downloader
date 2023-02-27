package testLowDisk

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/terminal"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	test_helpers.DirEmpty(progPath + "/Press-Stop/")

	progBounds := test_helpers.TestBounds(progPath)
	progBounds.LoadOption = consts.HIGH_LOAD // slow down so can stop after one file read
	globals.LogChannels = true
	misc.StartLog("../../../" + consts.CHANNEL_LOG_NAME)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.Press-Stop/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>
								<title>title tag</title>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-1.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-2.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-3.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-4.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-5.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-6.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-7.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-8.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-9.txt" length="17" type="text/plain" />
								</item>
								<item>
									<enclosure url="http://rss.Press-Stop/press-stop-10.txt" length="17" type="text/plain" />
								</item>
							</channel>
						</rss>`,
		"http://rss.Press-Stop/press-stop-1.txt":  `file 1 press-stop`,
		"http://rss.Press-Stop/press-stop-2.txt":  `file 2 press-stop`,
		"http://rss.Press-Stop/press-stop-3.txt":  `file 3 press-stop`,
		"http://rss.Press-Stop/press-stop-4.txt":  `file 4 press-stop`,
		"http://rss.Press-Stop/press-stop-5.txt":  `file 5 press-stop`,
		"http://rss.Press-Stop/press-stop-6.txt":  `file 6 press-stop`,
		"http://rss.Press-Stop/press-stop-7.txt":  `file 7 press-stop`,
		"http://rss.Press-Stop/press-stop-8.txt":  `file 8 press-stop`,
		"http://rss.Press-Stop/press-stop-9.txt":  `file 9 press-stop`,
		"http://rss.Press-Stop/press-stop-10.txt": `file 10 press-stop`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		httpResp := test_helpers.Http200Resp("rss.Press-Stop", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown mediaUrl : " + mediaUrl)
	return nil, nil
}

const expectedMenu string = `
1 |   0 files |    0MB | Press-Stop
 'Q' or a number + enter:
`

const expectedConsole string = `
Downloading 'Press-Stop' podcast, 10 files, hit 's' to stop
        TESTING - downloading stopped by simulated key press of 'S'
`

const expectedAdds = `
No changes detected
`

const expectedBads = `
`

func TestPressStopMock(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)

	globals.Console.Clear()

	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", actualMenu, err)
	}

	DurationOfTime := time.Duration(1) * time.Millisecond
	f := func() {
		keyStream <- "S"
	}
	time.AfterFunc(DurationOfTime, f)

	globals.Console.Clear()
	actualAdds, _ := terminal.AfterMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoice_1, httpTest)

	//timer1.Stop()
	//	fmt.Println("dddddddddddddddddd")
	//	fmt.Println("wa happen", actualAdds, podcastResults)
	//if !errors.Is(err, context.Canceled) {
	//if err.Error() != "TESTING - downloading stopped by simulated key press of 'S'" {
	//		fmt.Println("ddddddd  context.Canceled dddddd", err.Error())
	//	t.Fatal(err)
	//}

	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()
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
