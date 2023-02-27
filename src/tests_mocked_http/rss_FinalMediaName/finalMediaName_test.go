package testChooseName

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

/*
For BBC podcasts
  BBC News Top stories  go run ./ --emptyFiles podcasts.files.bbci.co.uk/p02nq0gn.rss
  Witness History       go run ./ --emptyFiles podcasts.files.bbci.co.uk/p004t1hd.rss

  Where :
    Content-Disposition
    attachment; filename="GlobalNewsPodcast-20230209-EarthquakeDeathsRiseToOver19000.mp3"


*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	test_helpers.DirRemove(progPath + "/final media name/")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

func httpTest(ctx context.Context, mediaUrl string) (*http.Response, error) {
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
									<enclosure url="http://rss.FinalMediaName/file-a.txt" length="11" type="text/plain" />
								</item>
							
								<item>
									<enclosure url="http://rss.FinalMediaName/file-b.txt" length="11" type="text/plain" />
								</item>

								<item>
								<title>My Title2</title>
									<enclosure url="http://rss.FinalMediaName/file-c.txt" length="11" type="text/plain" />
								</item>
							
								</channel>
						</rss>`,
		"http://rss.FinalMediaName/file-a.txt": `file a text`,
		"http://rss.FinalMediaName/file-b.txt": `file b text`,
		"http://rss.FinalMediaName/file-c.txt": `file c text`,
	}
	rssDisposition := map[string]string{
		"http://rss.FinalMediaName/podcast.xml": ``,
		"http://rss.FinalMediaName/file-a.txt":  `attachment; filename="final-a.txt"`,
		"http://rss.FinalMediaName/file-b.txt":  `attachment; filename="final-b.txt"`,
		"http://rss.FinalMediaName/file-c.txt":  `attachment; filename="final-c.txt"`,
	}

	if theData, ok := rssData[mediaUrl]; ok {
		thePath := rss.NameOfFile(mediaUrl)
		contentDisposition := rssDisposition[mediaUrl]
		httpResp := test_helpers.Http200Resp("rss.FinalMediaName", thePath, theData, contentDisposition)
		return httpResp, nil
	}
	fmt.Println("unknown FinalMediaName : " + mediaUrl)
	return nil, nil
}

const expectedConsole string = `
 Adding 'final media name'
        Downloading 'final media name' podcast, 3 files, hit 's' to stop
        	final-a.txt(read #0 11B)
        	final-b.txt(read #0 11B)
					final-c.txt(read #0 11B)
        		 final-a.txt (save #0, 0s)
        		 final-b.txt (save #0, 0s)
						 final-c.txt (save #0, 0s) 
`
const expectedAdds = `
Added 3 new files in 0s 
From http://rss.FinalMediaName/podcast.xml 
Into 'final media name'  
`
const expectedBads = ""

func Test_5_FinalMediaName(t *testing.T) {
	progBounds := setUp()
	cleanArgs := []string{"file-name.go", "http://rss.FinalMediaName/podcast.xml", "final", "media", "name"}
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, podcastResults := menu.ByNameOrUrl(cleanArgs, progBounds, keyStream, httpTest)

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

	podDir := progBounds.ProgPath + "/final media name"
	if _, err = os.Stat(podDir); err != nil {
		t.Fatal(podDir + " directory does not exist")
	}

	file1 := progBounds.ProgPath + "/final media name/final-a.txt"
	if _, err = os.Stat(file1); err != nil {
		t.Fatal(err)
	}

	file2 := progBounds.ProgPath + "/final media name/final-b.txt"
	if _, err = os.Stat(file2); err != nil {
		t.Fatal(err)
	}

	file3 := progBounds.ProgPath + "/final media name/final-c.txt"
	if _, err = os.Stat(file3); err != nil {
		t.Fatal(err)
	}

}
