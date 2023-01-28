package rss

//      go test ./...

//  const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-go-consol/main/src/tests/"

//                https://github.com/steenhansen/react-native-phone-recipes/blob/main/android/gradlew.bat
// https://raw.githubusercontent.com/steenhansen/react-native-phone-recipes/main/android/gradlew.bat

import (
	"errors"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
)

func TestInvalidXml(t *testing.T) {
	url := consts.TEST_DIR_URL + "invalid-xml/invalid-xml.rss"
	_, _, _, err := podcasts.ReadUrl(url)
	//  https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/tests/invalid-xml/invalid-xml.rss
	if !errors.Is(err, flaws.InvalidXML) {
		t.Fatalf(`TestInvalidXml failed`)
	}
}
