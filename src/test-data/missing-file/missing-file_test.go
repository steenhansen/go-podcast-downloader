package terminal

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/test-data/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

import (
	"errors"
	"fmt"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
)

func TestInvalidXml(t *testing.T) {
	pdescs := []string{"test pod desc"}
	feed := []string{consts.TEST_DIR_URL + "test-data/missing-file/missing-file.rss"}
	choice := 0
	progBounds := consts.ProgBounds{
		ProgPath:    "c:\\poddown",
		LoadOption:  "high",
		LimitOption: 0,
		MinDisk:     1000000000,
	}
	simKeyStream := make(chan string)
	mediaFix := map[string]error{}
	report, err := terminal.DownloadAndReport(pdescs, feed, choice, progBounds, simKeyStream, mediaFix)
	fmt.Println("test resort ", report)
	if !errors.Is(err, flaws.InvalidXML) {
		t.Fatalf(`TestInvalidXml failed`)
	}
}
