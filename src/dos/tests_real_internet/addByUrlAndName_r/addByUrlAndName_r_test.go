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

go run ./ siberiantimes.com/ecology/rss/ ecology

https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/tests_real_internet/AddByUrlAndName/git-server-source/add-by-url-and-name.rss

*/

func setUp() {
	progPath := misc.CurDir()
	testDir := progPath + "/add-by-url-and-name-r"
	test_helpers.DirRemove(testDir)

	//progBounds := test_helpers.TestBounds(progPath)
	//return progBounds
}

const expectedReport = `
Adding 'add-by-url-and-name-r'
        
        Downloading 'add-by-url-and-name-r' podcast, 3 files, hit 's' to stop
        	add-by-url-and-name-1.txt(read #0 21B)
        ERROR add-by-url-and-name-1.txt
        	add-by-url-and-name-2.txt(read #0 21B)
        ERROR add-by-url-and-name-2.txt
        	add-by-url-and-name-3.txt(read #0 21B)
        ERROR add-by-url-and-name-3.txt
`

func TestAddByUrlAndName(t *testing.T) {
	//progBounds := setUp()
	setUp()
	podcastUrl := consts.TEST_DIR_URL + "addByUrlAndName_r/git-server-source/add-by-url-and-name-r.rss"
	osArgs := []string{"AddByUrlAndName-test", podcastUrl, "add-by-url-and-name-r"}
	progBounds := test_helpers.TestBounds(misc.CurDir())
	keyStreamTest := make(chan string)
	_, podcastResults := terminal.AddByUrlAndName(podcastUrl, osArgs, progBounds, keyStreamTest, rss.HttpReal)
	if podcastResults.SeriousError != nil {
		t.Fatal(podcastResults.SeriousError)
	}
	actualReport := globals.Console.All()

	expectedDiff := test_helpers.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualReport), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedReport))
	}
}
