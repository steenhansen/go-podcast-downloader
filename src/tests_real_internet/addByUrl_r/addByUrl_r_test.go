package terminal

import (
	"testing"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/terminal"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

/*

go run ./ siberiantimes.com/ecology/rss/

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/addByUrl_r/git-server-source/add-by-url-r.rss

*/

func setUp() {
	progPath := misc.CurDir()
	globals.MediaMaxReadFileTime = consts.RSS_MAX_READ_FILE_TIME
	testDir := progPath + "/add-by-url-r"
	test_helpers.DirRemove(testDir)
	//progBounds := test_helpers.TestBounds(progPath)
	//return progBounds
}

const expectedReport = `
Adding 'add-by-url-r'
        
        Downloading 'add-by-url-r' podcast, 3 files, hit 's' to stop
        	add-by-url-1.txt(read #0 12B)
        		 add-by-url-1.txt (save #0, 0s)
        			Size disparity, expected 12 bytes, but was 12
        	add-by-url-2.txt(read #0 12B)
        		 add-by-url-2.txt (save #0, 0s)
        			Size disparity, expected 12 bytes, but was 12
        	add-by-url-3.txt(read #0 12B)
        		 add-by-url-3.txt (save #0, 0s)
        			Size disparity, expected 12 bytes, but was 12
`

func TestAddByUrl(t *testing.T) {
	//progBounds := setUp()
	setUp()
	podcastUrl := consts.TEST_DIR_URL + "addByUrl_r/git-server-source/add-by-url-r.rss"
	progBounds := test_helpers.TestBounds(misc.CurDir())
	keyStream := make(chan string)
	_, podcastResults := terminal.AddByUrl(podcastUrl, progBounds, keyStream, rss.HttpReal)
	if podcastResults.SeriousError != nil {
		t.Fatal(podcastResults.SeriousError)
	}
	//fmt.Println("res", res, err)
	actualReport := globals.Console.All()
	expectedDiff := test_helpers.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualReport), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedReport))
	}
}
