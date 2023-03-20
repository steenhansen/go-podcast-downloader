package misc

import (
	"regexp"
	"strconv"
	"strings"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
)

// go run ./console-downloader.go --minimumDisk=1_000_000_000_000_000
// will insist on there being 1 terrabyte of free disk space, or "low disk space" error
func SetMinDisk(osArgs []string) (int, []string, error) {
	theMinimum := consts.MIN_DISK_BYTES
	var err error
	minimumPlain := regexp.MustCompile(consts.MINIMUM_PLAIN)
	minimumAndNumber := regexp.MustCompile(consts.MINIMUM_AND_NUMBER)
	minimumNumber := regexp.MustCompile(consts.MINIMUM_NUMBER)
	limitArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && minimumPlain.MatchString(anArg) {
			if minimumAndNumber.MatchString(anArg) {
				limitStr := minimumNumber.FindString(anArg)
				noUnderscores := strings.ReplaceAll(limitStr, "_", "")
				theMinimum, err = strconv.Atoi(noUnderscores)
				if err != nil {
					return 0, nil, err
				}
			} else {
				return 0, nil, flaws.BadLimitSerious.MakeFlaw(anArg)
			}
		} else {
			limitArgs = append(limitArgs, anArg)
		}
	}
	return theMinimum, limitArgs, nil
}

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
