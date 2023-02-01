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

func ClampStr(testStr string) string {
	clampStr := "\n~" + testStr + "~"
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

func SameButOutOfOrder(realLines, expectedLines string) bool {
	reals := nonBlanks(realLines)
	expecteds := nonBlanks(expectedLines)
	for aReal := range reals {
		delete(expecteds, aReal)
	}
	return len(expecteds) == 0
}
