package terminal

import (
	"fmt"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

func setUp() consts.ProgBounds {
	progPath := misc.CurDir()
	delNotMissing := progPath + "/local-download-dest/not-missing.jpg"
	os.Remove(delNotMissing)
	progBounds := testings.ProgBounds(progPath)
	return progBounds
}

const expectedMenu string = ` 1 |              txt |   0 files |    1MB | local-download-dest
 'Q' or a number + enter: `
const expectedConsole string = `Downloading 'local-download-dest' podcast, hit 's' to stop
	no-such-file.txt(read #0)
	not-missing.txt(read #0)
	ERROR no-such-file.txt
	not-missing.txt (save #0, 0s)`
const expectedAdds = `
Added 1 new 'txt' file(s) in 0s 
From https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/missing-file.rss 
Into 'local-download-dest' `
const expectedBads = "\t\t*** 404 or 400 html page, https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/git-server-source/no-such-file.txt\n"

/// just use .txt files

func TestMissingFileFromMenu(t *testing.T) {
	progBounds := setUp()
	simKeyStream := make(chan string)
	globals.Console.Clear()
	podcastMenu, err := terminal.ShowNumberedChoices(progBounds)
	if err != nil {
		fmt.Println("wa happen", err)
	}
	globals.Console.Clear()
	addedFiles, _ := terminal.AfterMenu(progBounds, simKeyStream, testings.KeyboardMenuChoice_1)
	actualConsol := globals.Console.All()
	badFiles := globals.Faults.All()
	if podcastMenu != expectedMenu {
		t.Fatal(testings.ClampActual(podcastMenu), testings.ClampExpected(expectedMenu))
	}

	if !testings.SameButOutOfOrder(actualConsol, expectedConsole) {
		t.Fatal(testings.ClampActual(actualConsol), testings.ClampExpected(expectedConsole))
	}

	if expectedAdds != addedFiles {
		t.Fatal(testings.ClampActual(addedFiles), testings.ClampExpected(expectedAdds))
	}
	if expectedBads != badFiles {
		t.Fatal(testings.ClampActual(badFiles), testings.ClampExpected(expectedBads))
	}
}
