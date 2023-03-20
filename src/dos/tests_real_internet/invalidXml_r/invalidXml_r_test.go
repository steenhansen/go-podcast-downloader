package rss

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/dos/tests_real_internet/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

import (
	"errors"
	"testing"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/podcasts"
	"podcast-downloader/src/dos/rss"
)

func TestInvalidXml_r(t *testing.T) {
	url := consts.TEST_DIR_URL + "invalidXml_r/invalid-xml-r.rss"

	_, _, _, _, err := podcasts.ReadRssUrl(url, rss.HttpReal)

	//  https://raw.githubusercontent.com/steenhansen/go-podcast-downloader/main/src/dos/tests_real_internet/invalid-xml/invalid-xml.rss
	if !errors.Is(err, flaws.InvalidXML) {
		t.Fatal(err)
	}
}
