package testLowDisk

import (
	"errors"
	"testing"

	"podcast-downloader/src/dos/console"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/test_helpers"
)

/*
https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/lowDisk_r/git-server-source/low-disk-r.rss
*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	progBounds := test_helpers.TestBounds(progPath)
	progBounds.MinDisk = 1_000_000_000_000_000
	return progBounds
}

const expectedConsole string = `
 1 |   0 files |    0MB | low-disk-r
         'Q' or a number + enter: Downloading 'low-disk-r' podcast, 2 files, hit 's' to stop
        	low-disk-r-1.txt(read #0 11B)
        ERROR: E_15
        low disk space, 22GB free, need minimum 909TB to proceed FILE: low-disk-r-1.txt
`

const expectedAdds = `
No changes detected
`

const expectedBads = `
E_15 : low disk space, 96GB free, need minimum 909TB to proceed
`

func TestLowDisk_r(t *testing.T) {
	progBounds := setUp()
	keyStreamTest := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := console.DisplayMenu(progBounds, keyStreamTest, test_helpers.KeyboardMenuChoiceNum("1"), rss.HttpReal)
	var flawError flaws.FlawError
	err := podcastResults.SeriousError
	if errors.As(err, &flawError) {
		lowErr := flawError.Error()
		safeErr := test_helpers.ReplaceXxGbFree(lowErr)
		if safeErr != "low disk space, xxGB free, need minimum 909TB to proceed" {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	actualConsole := globals.Console.All()
	actualBads := globals.Faults.All()

	expectedDiff := test_helpers.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualConsole), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedConsole))
	}

	if test_helpers.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(test_helpers.ClampActual(actualAdds), test_helpers.ClampExpected(expectedAdds))
	}

	actualSafe := test_helpers.ReplaceXxGbFree(actualBads)
	expectedSafe := test_helpers.ReplaceXxGbFree(expectedBads)
	if test_helpers.NotSameTrimmed(actualSafe, expectedSafe) {
		t.Fatal(test_helpers.ClampActual(actualSafe), test_helpers.ClampExpected(expectedSafe))
	}

}
