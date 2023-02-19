package misc

import (
	"regexp"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
)

func setDnsErrors(osArgs []string) []string {
	dnsErrors := regexp.MustCompile(consts.DNS_ERRORS)
	emptyArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && dnsErrors.MatchString(anArg) {
			globals.DnsErrorsTest = true
		} else {
			emptyArgs = append(emptyArgs, anArg)
		}
	}
	return emptyArgs
}

func setEmptyFiles(osArgs []string) []string {
	emptyFiles := regexp.MustCompile(consts.EMTPY_FILES)
	emptyArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && emptyFiles.MatchString(anArg) {
			globals.EmptyFilesTest = true
		} else {
			emptyArgs = append(emptyArgs, anArg)
		}
	}
	return emptyArgs
}

func IsTesting(osArgs []string) bool {
	for _, anArg := range osArgs {
		if strings.HasPrefix(anArg, consts.TEST_FLAG_PREFIX) {
			return true
		}
	}
	return false
}

// go run ./
// go run ./  -race
func DelRace(osArgs []string) ([]string, error) {
	raceArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex == 0 || anArg != consts.RACE_DEBUG {
			raceArgs = append(raceArgs, anArg)
		}
	}
	return raceArgs, nil
}
