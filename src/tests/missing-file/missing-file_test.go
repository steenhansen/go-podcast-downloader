package terminal

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/tests/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

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
	delNotMissing := progPath + "/local-dest/not-missing.jpg"
	os.Remove(delNotMissing)
	progBounds := testings.ProgBounds(progPath)
	return progBounds
}

const expectedMenu string = ` 1 |                  |   0 files |    0MB | local-dest
 'Q' or a number + enter: `
const expectedConsole string = `Downloading 'local-dest' podcast, hit 's' to stop
	no-such-file.jpg(read #0)
	not-missing.jpg(read #0)
	ERROR no-such-file.jpg
	not-missing.jpg (save #0, 0s)`
const expectedAdds = `
Added 1 new 'jpg' file(s) in 0s 
From https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/server-source/missing-file.rss 
Into 'local-dest' `
const expectedBads = "\t\t*** 404 or 400 html page, https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/server-source/no-such-file.jpg\n"

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
		t.Fatal(testings.ClampStr(podcastMenu), testings.ClampStr(expectedMenu))
	}

	if !testings.SameButOutOfOrder(actualConsol, expectedConsole) {
		t.Fatal(testings.ClampStr(actualConsol), testings.ClampStr(expectedConsole))
	}

	if expectedAdds != addedFiles {
		t.Fatal(testings.ClampStr(addedFiles), testings.ClampStr(expectedAdds))
	}
	if expectedBads != badFiles {
		t.Fatal(testings.ClampStr(badFiles), testings.ClampStr(expectedBads))
	}
}
