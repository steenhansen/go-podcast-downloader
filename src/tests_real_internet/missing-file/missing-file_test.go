package terminal

import (
	"fmt"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/test_helpers"
)

/*

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/missing-file/git-server-source/missing-file.rss

*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	os.Remove(progPath + "/local-download-dest/not-missing.txt")
	progBounds := test_helpers.TestBounds(progPath)
	return progBounds
}

const expectedMenu string = `
1 |   0 files |    0MB | local-download-dest
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
Added 1 new files in 0s
From https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/missing-file/git-server-source/missing-file.rss
Into 'local-download-dest'
`

const expectedBads = `	
E_10 : HTTP error 404 Not Found : https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/missing-file/git-server-source/no-such-file.txt
`

func TestMissingFileFromMenu(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	globals.Console.Clear()
	actualAdds, err := terminal.AfterMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoice_1, rss.HttpReal)
	if err != nil {
		t.Fatal(err)
	}
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
