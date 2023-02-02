package testings

import (
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
	return ClampStr(testStr, "actual")
}

func ClampExpected(testStr string) string {
	return ClampStr(testStr, "expected")
}

func ClampStr(testStr string, actualStr string) string {
	clampStr := "\n"
	if actualStr == "actual" {
		clampStr += "ACTUAL:"
	} else {
		clampStr += "EXPECT:"

	}
	clampStr += "~" + testStr + "~\n"
	return clampStr
}

func KeyboardMenuChoice_1() string {
	return "1"
}

func ProgBounds(progPath string) consts.ProgBounds {
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
		if aLine != "" {
			trimmedLine := strings.TrimSpace(aLine)
			textLines[trimmedLine] = trimmedLine
		}
	}
	return textLines
}

func NotSameOutOfOrder(actualLines, expectedLines string) bool {
	trimmedActual := strings.TrimSpace(actualLines)
	trimmedExpected := strings.TrimSpace(expectedLines)
	actuals := nonBlanks(trimmedActual)
	expecteds := nonBlanks(trimmedExpected)
	for aLine := range actuals {
		delete(actuals, aLine)
		delete(expecteds, aLine)
	}
	//	fmt.Println("actuals", actuals)
	//fmt.Println("expecteds", expecteds)
	return len(expecteds) != 0
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
