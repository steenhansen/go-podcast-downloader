package testTimeout

import (
	"errors"
	"testing"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/menu"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	progBounds := test_helpers.TestBounds(progPath)
	progBounds.LoadOption = consts.LOW_LOAD
	globals.MediaMaxReadFileTime = time.Microsecond * 1
	return progBounds
}

const expectedConsole string = `
1 |  10 files |    0MB | timeout-r
         'Q' or a number + enter: Downloading 'timeout-r' podcast, 10 files, hit 's' to stop
        				Have #1 file-1.txt
        				Have #2 file-2.txt
        				Have #3 file-3.txt
        				Have #4 file-4.txt
        				Have #5 file-5.txt
        				Have #6 file-6.txt
        				Have #7 file-7.txt
        				Have #8 file-8.txt
        				Have #9 file-9.txt
        				Have #10 file-10.txt
`
const expectedAdds = `
No changes detected
`

func TestTimeout_m(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)
	globals.Console.Clear()
	actualAdds, _, podcastResults := menu.DisplayMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoiceNum("1"), rss.HttpReal)
	var flawError flaws.FlawError
	err := podcastResults.SeriousError
	if errors.As(err, &flawError) {
		timeoutErr := flawError.Error()
		if timeoutErr != "Internet connection timed out by exceeding duration: 1µs" {
			t.Fatal(err)
		}
	} else {
		t.Fatal(err)
	}
	actualConsole := globals.Console.All()
	expectedDiff := test_helpers.NotSameOutOfOrder(actualConsole, expectedConsole)
	if len(expectedDiff) != 0 {
		t.Fatal(test_helpers.ClampActual(actualConsole), test_helpers.ClampMapDiff(expectedDiff), test_helpers.ClampExpected(expectedConsole))
	}

	if test_helpers.NotSameTrimmed(actualAdds, expectedAdds) {
		t.Fatal(test_helpers.ClampActual(actualAdds), test_helpers.ClampExpected(expectedAdds))
	}

}
