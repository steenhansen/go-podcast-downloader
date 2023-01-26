package rss

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/test-data/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
)

func TestInvalidXml(t *testing.T) {
	url := consts.TEST_DIR_URL + "invalid-xml/invalid-xml.rss"
	_, _, _, err := podcasts.ReadUrl(url)
	fmt.Println("erddddddddddddddddddddddddddddr err", err)
	fmt.Println("erddddddddddddddddddddddddddddr url", url)
	fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXX args", os.Args)
	// C:\Users\16043\AppData\Local\Temp\go-build897778877\b001\invalid-xml.test.exe                        // test
	//args [C:\Users\16043\AppData\Local\Temp\go-build4259596642\b001\invalid-xml.test.exe
	//-test.testlogfile=C:\Users\16043\AppData\Local\Temp\go-build4259596642\b001\testlog.txt
	//-test.paniconexit0 -test.timeout=30s -test.run=^TestInvalidXml$]
	// C:\Users\16043\AppData\Local\Temp\go-build4070268852\b001\exe\go-podcast-downloader-console.exe                   // comand line
	if errors.Is(err, flaws.InvalidXML) {
		t.Fatalf(`TestInvalidXml failed`)
	}
}
