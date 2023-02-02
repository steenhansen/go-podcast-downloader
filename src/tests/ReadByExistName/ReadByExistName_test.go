package terminal

import (
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

/*

go run ./ siberiantimes.com/ecology/rss/ ecology

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/ReadByExistName/git-server-source/read-by-exist-name.rss

*/

func setUp() {
	progPath := misc.CurDir()
	testDir := progPath + "/local-download-dest"
	testings.DirRemove(testDir)
}

const expectedReport = `
Adding 'local-download-dest'
Downloading 'local-download-dest' podcast, hit 's' to stop
read-by-exist-name-1.txt(read #0 21B)
read-by-exist-name-2.txt(read #0 21B)
read-by-exist-name-3.txt(read #0 21B)
read-by-exist-name-1.txt (save #0, 0s)
read-by-exist-name-2.txt (save #0, 0s)
read-by-exist-name-3.txt (save #0, 0s)
`

func TestReadByExistName(t *testing.T) {
	setUp()
	podcastUrl := consts.TEST_DIR_URL + "ReadByExistName/git-server-source/read-by-exist-name.rss"
	osArgs := []string{"ReadByExistName-test", podcastUrl, "local-download-dest"}
	progBounds := testings.ProgBounds(misc.CurDir())
	simKeyStream := make(chan string)
	terminal.ReadByExistName(osArgs, progBounds, simKeyStream)
	actualReport := globals.Console.All()
	if testings.NotSameOutOfOrder(actualReport, expectedReport) {
		t.Fatal(testings.ClampActual(actualReport), testings.ClampExpected(expectedReport))
	}
}
