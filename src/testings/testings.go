package testings

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

func DirRemove(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		return err
	}
	err = os.Remove(dirPath)
	return err
}

func ClampActual(testStr string) string {
	return ClampStr(testStr, "ACTUAL : ")
}

func ClampExpected(testStr string) string {
	return ClampStr(testStr, "EXPECTED : ")
}

func ClampMapDiff(expectedDiff map[string]string) (mapDiff string) {
	for v := range expectedDiff {
		mapDiff += "\n" + v
	}
	return ClampStr(mapDiff, "DIFFERENCE : ")
}

func ClampStr(testStr string, actualStr string) string {
	clampStr := "\n" + actualStr + "~~~" + testStr + "~~~\n"
	return clampStr
}

func KeyboardMenuChoice_1() string {
	return "1"
}

// testings.KeyboardMenuChoiceNum("q")
// testings.KeyboardMenuChoiceNum("12") always choose "12" on menu
func KeyboardMenuChoiceNum(simChoice string) func() string {
	menuChoice := func() string {
		return simChoice
	}
	return menuChoice
}

func TestBounds(progPath string) consts.ProgBounds {
	progBounds := consts.ProgBounds{
		ProgPath:    progPath,
		LoadOption:  consts.HIGH_LOAD,
		LimitOption: 0,
		MinDisk:     1000000000,
	}
	return progBounds
}

func nonBlanks(consoleOutput string) map[string]string {
	textLines := make(map[string]string)
	consoleLines := strings.Split(consoleOutput, "\n")
	for _, aLine := range consoleLines {
		trimmedLine := strings.TrimSpace(aLine)
		if trimmedLine != "" {
			textLines[trimmedLine] = trimmedLine
		}
	}
	return textLines
}

func NotSameOutOfOrder(actualLines, expectedLines string) map[string]string {
	trimmedActual := strings.TrimSpace(actualLines)
	trimmedExpected := strings.TrimSpace(expectedLines)
	actuals := nonBlanks(trimmedActual)
	expecteds := nonBlanks(trimmedExpected)
	for aLine := range actuals {
		delete(actuals, aLine)
		delete(expecteds, aLine)
	}
	return expecteds
}

func NotSameTrimmed(actualStr, expectedStr string) bool {
	trimmedActual := strings.TrimSpace(actualStr)
	trimmedExpected := strings.TrimSpace(expectedStr)
	actualLines := strings.Split(trimmedActual, "\n")
	expectedLines := strings.Split(trimmedExpected, "\n")
	for i, actual := range actualLines {
		actualLines[i] = strings.TrimSpace(actual)
	}
	for i, expected := range expectedLines {
		expectedLines[i] = strings.TrimSpace(expected)
	}
	cleanActual := strings.Join(actualLines, "\n")
	cleanExpected := strings.Join(expectedLines, "\n")
	return cleanActual != cleanExpected
}

// https://stackoverflow.com/questions/33978216/create-http-response-instance-with-sample-body-string-in-golang
func Http200Resp(theHost, thePath, bodyXml string) *http.Response {

	theUrl := &url.URL{
		Scheme: "http",
		Host:   theHost,
		Path:   thePath,
	}

	httpReq := &http.Request{
		URL: theUrl,
	}

	httpResp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(bodyXml)),
		ContentLength: int64(len(bodyXml)),
		Request:       httpReq,
		Header:        make(http.Header, 0),
	}
	return httpResp
}
