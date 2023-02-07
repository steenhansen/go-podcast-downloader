package terminal

import (
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

/*

go run ./ siberiantimes.com/ecology/rss/

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/internet-tests/AddByUrl/git-server-source/add-by-url.rss

*/

func setUp() {
	progPath := misc.CurDir()
	testDir := progPath + "/add-by-url-title-97"
	testings.DirRemove(testDir)
	//progBounds := testings.TestBounds(progPath)
	//return progBounds
}

const expectedReport = `
Adding 'add-by-url-title-97'
Downloading 'add-by-url-title-97' podcast, 3 files, hit 's' to stop
add-by-url-1.txt(read #0 21B)
add-by-url-2.txt(read #0 21B)
add-by-url-3.txt(read #0 21B)
add-by-url-1.txt (save #0, 0s)
add-by-url-2.txt (save #0, 0s)
add-by-url-3.txt (save #0, 0s)
`

func TestAddByUrl(t *testing.T) {
	//progBounds := setUp()
	setUp()
	podcastUrl := consts.TEST_DIR_URL + "AddByUrl/git-server-source/add-by-url.rss"
	progBounds := testings.TestBounds(misc.CurDir())
	keyStream := make(chan string)
	terminal.AddByUrl(podcastUrl, progBounds, keyStream, rss.HttpMedia)
	actualReport := globals.Console.All()
	expectedDiff := testings.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(testings.ClampActual(actualReport), testings.ClampMapDiff(expectedDiff), testings.ClampExpected(expectedReport))
	}
}
