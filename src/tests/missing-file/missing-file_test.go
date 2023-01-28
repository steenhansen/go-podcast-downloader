package terminal

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/tests/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

import (
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/tests/testFuncs"
)

func TestMissingFileFromMenu(t *testing.T) {

	path := media.CurDir()
	delNotMissing := path + "/local-dest/not-missing.jpg"
	e := os.Remove(delNotMissing)
	if e != nil {
		t.Fatal(e)
	}
	progBounds := misc.TestProgBounds(path)
	//	simKeyStream := make(chan string)
	misc.ConsolOutput = ""
	theMenu, _ := terminal.ShowNumberedChoices(progBounds)
	misc.ConsolOutput = ""
	//	report, _ := terminal.AfterMenu(progBounds, simKeyStream, misc.GetMenuChoiceTest1)
	expectedMenu := ` 1 |              |   0 files |    0MB | test pod desc
 'Q' or a number + enter: `
	// 	expectedFiles := `Downloading 'test pod desc' podcast, hit 's' to stop
	// no-such-file.jpeg(read #0)
	// not-missing.jpeg(read #0)
	// ERROR no-such-file.jpeg
	// not-missing.jpeg (save #0, 0s)`

	// 	expectedAddition := `Added 0 new 'jpeg' file(s) from https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/missing-file.rss into 'test pod desc' in 0s`

	// 	expectedError := "\t\t*** 404 or 400 html page, https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/missing-file/no-such-file.jpeg\n"

	// 	badFiles := misc.GetMediaFaults2()

	if theMenu != expectedMenu {
		t.Fatal(testFuncs.ClampStr(theMenu), testFuncs.ClampStr(expectedMenu))
	}
	// if misc.SameButOutOfOrder(misc.ConsolOutput, expectedFiles) {
	// 	t.Fatal(misc.ClampStr(misc.ConsolOutput), misc.ClampStr(expectedFiles))
	// }
	// if expectedAddition != report {
	// 	t.Fatal(misc.ClampStr(expectedAddition), misc.ClampStr(report))
	// }
	// if expectedError != badFiles {
	// 	t.Fatal(misc.ClampStr(expectedError), misc.ClampStr(badFiles))
	// }
}
