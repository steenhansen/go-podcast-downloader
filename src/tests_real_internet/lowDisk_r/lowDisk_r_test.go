package testLowDisk

import (
	"errors"
	"fmt"
	"testing"

	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

/*

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/lowDisk_r/git-server-source/low-disk-r.rss

*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	progBounds := test_helpers.TestBounds(progPath)
	progBounds.MinDisk = 1_000_000_000_000_000
	return progBounds
}

const expectedConsole string = `
 1 |   1 files |    0MB | Low-Disk
 'Q' or a number + enter: Downloading 'low-disk-r' podcast, 2 files, hit 's' to stop
				Have #1 low-disk-r-1.txt
	low-disk-r-1.txt(read #0 16B)
ERROR low-disk-r-2.txt
`
const expectedAdds = `
No changes detected
`

const expectedBads = `
E_15 : low disk space, 96GB free, need minimum 909TB to proceed
`

func TestLowDisk_m(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := menu.DisplayMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoiceNum("1"), rss.HttpReal)
	var flawError flaws.FlawError
	err := podcastResults.SeriousError
	if errors.As(err, &flawError) {
		if flawError.Error() != "low disk space, 96GB free, need minimum 909TB to proceed" {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()
	fmt.Println("------------")

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
