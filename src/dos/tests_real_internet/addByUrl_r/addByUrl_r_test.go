package terminal

import (
	"testing"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/terminal"
	"podcast-downloader/src/dos/test_helpers"
)

/*

go run ./ siberiantimes.com/ecology/rss/

https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/addByUrl_r/git-server-source/add-by-url-r.rss

*/

func setUp() {
	progPath := misc.CurDir()
	testDir := progPath + "/add-by-url-r"
	test_helpers.DirRemove(testDir)
	//progBounds := test_helpers.TestBounds(progPath)
	//return progBounds
}

const expectedReport = `
       Adding 'add-by-url-r'

        Downloading 'add-by-url-r' podcast, 3 files, hit 's' to stop
                add-by-url-1.txt(read #0 12B)
                add-by-url-2.txt(read #0 12B)
                         add-by-url-1.txt (save #0, 0s)
                                Size disparity, expected 12 bytes, but was 12
                         add-by-url-2.txt (save #0, 0s)
                                Size disparity, expected 12 bytes, but was 12
                add-by-url-3.txt(read #0 12B)
                         add-by-url-3.txt (save #0, 0s)
                                Size disparity, expected 12 bytes, but was 12
`

func TestAddByUrl(t *testing.T) {
	setUp()
	podcastUrl := consts.TEST_DIR_URL + "addByUrl_r/git-server-source/add-by-url-r.rss"
	progBounds := test_helpers.TestBounds(misc.CurDir())
	keyStreamTest := make(chan string)
	_, podcastResults := terminal.AddByUrl(podcastUrl, progBounds, keyStreamTest, rss.HttpReal)
	if podcastResults.SeriousError != nil {
		t.Fatal(podcastResults.SeriousError)
	}
	actualReport := globals.Console.All()
	expectedDiff := test_helpers.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualReport), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedReport))
	}
}
