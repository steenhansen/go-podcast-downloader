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

go run ./ ecology

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/ReadByExistName/git-server-source/read-by-exist-name.rss

*/

const expectedReport = `
Downloading 'read-by-exist-name-r' podcast, 3 files, hit 's' to stop
        				Have #1 read-by-exist-name-1.txt
        				Have #2 read-by-exist-name-2.txt
        				Have #3 read-by-exist-name-3.txt
`

func TestReadByExistName_r(t *testing.T) {
	podcastUrl := consts.TEST_DIR_URL + "ReadByExistName_r/git-server-source/read-by-exist-name-r.rss"
	osArgs := []string{"ReadByExistName-test", podcastUrl, "read-by-exist-name-r"}
	progBounds := test_helpers.TestBounds(misc.CurDir())
	keyStream := make(chan string)
	_, podcastResults := terminal.ReadByExistName(osArgs, progBounds, keyStream, rss.HttpReal)
	if podcastResults.SeriousError != nil {
		t.Fatal(podcastResults.SeriousError)
	}
	actualReport := globals.Console.All()

	expectedDiff := test_helpers.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualReport), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedReport))
	}

}
