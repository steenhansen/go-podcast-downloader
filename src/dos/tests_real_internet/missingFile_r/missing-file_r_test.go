package t1

import (
	"fmt"
	"os"
	"testing"
	"time"

	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/terminal"
	"podcast-downloader/src/dos/test_helpers"
)

/*

https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/missingFile_r/git-server-source/missing-file-r.rss

*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()

	os.Remove(progPath + "/missing-file-r/not-missing.txt")
	progBounds := test_helpers.TestBounds(progPath)

	globals.MediaMaxReadFileTime = time.Second * 5

	return progBounds
}

const expectedMenu string = `
1 |   0 files |    0MB | missing-file-r
 'Q' or a number + enter:
`

const expectedConsole string = `
Downloading 'missing-file-r' podcast, 2 files, hit 's' to stop
                not-missing.txt(read #0 11B)
                         not-missing.txt (save #0, 0s)
                                Size disparity, expected 11 bytes, but was 11
                no-such-file.txt(read #0 12B)
        ERROR: E_10
        HTTP error 404 Not Found : https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/missingFile_r/git-server-source/no-such-file.txt FILE: no-such-file.txt
`

const expectedAdds = `
Added 1 new files in 0s
From https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/missingFile_r/git-server-source/missing-file-r.rss
Into 'missing-file-r'
`

const expectedBads = `	
E_10 : HTTP error 404 Not Found : https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/missingFile_r/git-server-source/no-such-file.txt
`

func TestMissingFile_r(t *testing.T) {
	progBounds := setUp()
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	globals.Console.Clear()
	getMenuChoice := test_helpers.KeyboardMenuChoice_1
	actualAdds, _ := terminal.AfterMenu(progBounds, keyStreamTest, getMenuChoice, rss.HttpReal)
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()
	if test_helpers.NotSameTrimmed(actualMenu, expectedMenu) {
		t.Fatal(test_helpers.ClampActual(actualMenu), test_helpers.ClampExpected(expectedMenu))
	}

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
}
