package t2

import (
	"fmt"
	"testing"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/terminal"
	"github.com/steenhansen/go-podcast-downloader/src/test_helpers"
)

/*

https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/press-stop/git-server-source/press-stop.rss

*/

func setUp() models.ProgBounds {
	progPath := misc.CurDir()
	test_helpers.DirEmpty(progPath + "/title tag/")

	progBounds := test_helpers.TestBounds(progPath)
	progBounds.LoadOption = consts.LOW_LOAD // slow down so can stop after one file read
	globals.LogChannels = true
	misc.StartLog("/test-chan-log.txt")
	return progBounds
}

const expectedMenu string = `
1 |   0 files |    0MB | title tag
 'Q' or a number + enter:
`

const expectedConsole string = `
Downloading 'title tag' podcast, 10 files, hit 's' to stop
        TESTING - downloading stopped by simulated key press of 'S'
        	press-stop-1.txt(read #0 12B)
        	press-stop-2.txt(read #0 12B)
`

const expectedAdds = `
No changes detected
`

//E_10 : HTTP error 404 Not Found : https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests_real_internet/press-stop/git-server-source/no-such-file.txt

const expectedBads = `
`

//    go test ./src/tests_real_internet/press-stop/... -count=1 -timeout 22s           OK

func TestMissingFileFromMenu(t *testing.T) {
	progBounds := setUp()
	keyStream := make(chan string)

	globals.Console.Clear()

	actualMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", actualMenu, err)
	}

	DurationOfTime := time.Duration(1) * time.Second
	f := func() {
		fmt.Println("aaaaaaaaaaaaaaa")
		keyStream <- "S"
	}
	time.AfterFunc(DurationOfTime, f)

	globals.Console.Clear()
	fmt.Println("bbbbbbbbbbbbbbbb")
	actualAdds, podcastResults := terminal.AfterMenu(progBounds, keyStream, test_helpers.KeyboardMenuChoice_1, rss.HttpReal)
	fmt.Println("ccccccccccccccccccccc")

	//timer1.Stop()
	fmt.Println("dddddddddddddddddd")
	fmt.Println("wa happen", actualAdds, podcastResults)
	//if !errors.Is(err, context.Canceled) {
	//if err.Error() != "TESTING - downloading stopped by simulated key press of 'S'" {
	//		fmt.Println("ddddddd  context.Canceled dddddd", err.Error())
	//	t.Fatal(err)
	//}

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
