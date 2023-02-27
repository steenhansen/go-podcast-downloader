package testChooseName

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

/*
  For feeds that have every podcast episode named the same

	Like Stuff You Should Know,

     go run ./ --emptyFiles omnycontent.com/d/playlist/e73c998e-6e60-432f-8610-ae210140c5b1/A91018A4-EA4F-4130-BF55-AE270180C327/44710ECC-10BB-48D1-93C7-AE270180C33E/podcast.rss


		 https://chtbl.com/track/5899E/podtrac.com/pts/redirect.mp3/traffic.omny.fm/d/clips/e7 ... f12/audio.mp3
*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	globals.MediaMaxReadFileTime = consts.RSS_MAX_READ_FILE_TIME
	test_helpers.DirRemove(progPath + "/choose-name-m/")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
	if ctx.Err() == context.Canceled {
		return nil, context.Canceled
	}
	rssData := map[string]string{
		"http://rss.chooseName/podcast.xml": `<?xml version="1.0" encoding="UTF-8"?>
						<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:atom="http://www.w3.org/2005/Atom">
							<channel>

								<title>title tag</title>
								<itunes:title>changed tag</itunes:title>
								<item>
								<title>first-title</title>
									<enclosure url="http://rss.chooseName/re-direct.txt?finFileName=first.text" length="21" type="text/plain" />
								</item>
							
								<item>
							<title>second-title</title>
									<enclosure url="http://rss.chooseName/re-direct.txt?finFileName=second.text" length="11" type="text/plain" />
								</item>

								<item>
								<title>third-title</title>
									<enclosure url="http://rss.chooseName/re-direct.txt?finFileName=third.text" length="33" type="text/plain" />
								</item>
							
								</channel>
						</rss>`,
		"http://rss.chooseName/re-direct.txt?finFileName=first.text":  `first file with title`,
		"http://rss.chooseName/re-direct.txt?finFileName=second.text": `second file`,
		"http://rss.chooseName/re-direct.txt?finFileName=third.text":  `third file also has a valid title`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := ""
		//fmt.Println("mediaUrl=", mediaUrl, "++++ thePath=", thePath, "::: theData=", theData)
		httpResp := test_helpers.Http200Resp("rss.chooseName", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown chooseName : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `
Adding 'choose-name-m'
        
        Downloading 'choose-name-m' podcast, 3 files, hit 's' to stop
        	first-title.txt(read #0 21B)
        	third-title.txt(read #0 33B)
        	second-title.txt(read #0 11B)
        		 first-title.txt (save #0, 0s)
        			Size disparity, expected 21 bytes, but was 21
        		 third-title.txt (save #0, 0s)
        			Size disparity, expected 33 bytes, but was 33
        		 second-title.txt (save #0, 0s)
        			Size disparity, expected 11 bytes, but was 11
`
const expectedAdds = `
Added 3 new files in 0s 
From http://rss.chooseName/podcast.xml 
Into 'choose-name-m'  
`
const expectedBads = ""

func Test_1_ByNameOrUrl_AddByUrlAndName(t *testing.T) {
	progBounds := setUp()
	cleanArgs := []string{"file-name.go", "http://rss.chooseName/podcast.xml", "choose-name-m"}
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, podcastResults := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, httpTest)
	//fmt.Println("actualAdds", actualAdds)
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

	podDir := progBounds.ProgPath + "/choose-name-m"
	if _, err = os.Stat(podDir); err != nil {
		t.Fatal(podDir + " directory does not exist")
	}

	file1 := progBounds.ProgPath + "/choose-name-m/first-title.txt"
	if _, err = os.Stat(file1); err != nil {
		t.Fatal(err)
	}

	file2 := progBounds.ProgPath + "/choose-name-m/second-title.txt"
	if _, err = os.Stat(file2); err != nil {
		t.Fatal(err)
	}

	file3 := progBounds.ProgPath + "/choose-name-m/third-title.txt"
	if _, err = os.Stat(file3); err != nil {
		t.Fatal(err)
	}

}
