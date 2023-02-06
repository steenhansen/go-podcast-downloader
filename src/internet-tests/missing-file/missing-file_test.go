package terminal

import (
	"fmt"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

/*

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/missing-file.rss

*/

func setUp() consts.ProgBounds {
	progPath := misc.CurDir()
	os.Remove(progPath + "/local-download-dest/not-missing.txt")
	progBounds := testings.TestBounds(progPath)
	return progBounds
}

const expectedMenu string = `
1 |                  |   0 files |    0MB | local-download-dest
 'Q' or a number + enter:
`

const expectedConsole string = `
Downloading 'local-download-dest' podcast, 2 files, hit 's' to stop
no-such-file.txt(read #0 12B)
not-missing.txt(read #0 11B)
ERROR no-such-file.txt
not-missing.txt (save #0, 0s)
`

const expectedAdds = `
Added 1 new 'txt' file(s) in 0s
From https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/missing-file.rss
Into 'local-download-dest'
`

const expectedBads = "\t\t*** 404 or 400 html page, https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/no-such-file.txt\n"

func TestMissingFileFromMenu(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	globals.Console.Clear()
	actualAdds, _ := terminal.AfterMenu(progBounds, keyStream, testings.KeyboardMenuChoice_1, rss.HttpMedia)
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()
	if testings.NotSameTrimmed(actualMenu, expectedMenu) {
		t.Fatal(testings.ClampActual(actualMenu), testings.ClampExpected(expectedMenu))
	}

	expectedDiff := testings.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(testings.ClampActual(actualConsole), testings.ClampMapDiff(expectedDiff), testings.ClampExpected(expectedConsole))
	}

	if testings.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(testings.ClampActual(actualAdds), testings.ClampExpected(expectedAdds))
	}

	if testings.NotSameTrimmed(actualBads, expectedBads) {
		t.Fatal(testings.ClampActual(actualBads), testings.ClampExpected(expectedBads))
	}
}
