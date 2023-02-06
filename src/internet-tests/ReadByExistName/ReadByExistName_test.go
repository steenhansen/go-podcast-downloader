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

go run ./ ecology

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/ReadByExistName/git-server-source/read-by-exist-name.rss

*/

const expectedReport = `
Downloading 'local-download-dest' podcast, 3 files, hit 's' to stop
`

func TestReadByExistName(t *testing.T) {
	podcastUrl := consts.TEST_DIR_URL + "ReadByExistName/git-server-source/read-by-exist-name.rss"
	osArgs := []string{"ReadByExistName-test", podcastUrl, "local-download-dest"}
	progBounds := testings.TestBounds(misc.CurDir())
	keyStream := make(chan string)
	terminal.ReadByExistName(osArgs, progBounds, keyStream, rss.HttpMedia)
	actualReport := globals.Console.All()

	expectedDiff := testings.NotSameOutOfOrder(actualReport, expectedReport)
	if len(expectedDiff) != 0 {
		t.Fatal(testings.ClampActual(actualReport), testings.ClampMapDiff(expectedDiff), testings.ClampExpected(expectedReport))
	}

}
