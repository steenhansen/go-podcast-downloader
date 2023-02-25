package misc

import (
	"regexp"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
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

func setLogChannels(osArgs []string) []string {
	logChannels := regexp.MustCompile(consts.LOG_CHANNELS)
	emptyArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && logChannels.MatchString(anArg) {
			globals.LogChannels = true
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
