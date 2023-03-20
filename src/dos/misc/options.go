package misc

import (
	"fmt"
	"regexp"
	"strconv"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
)

func setForceTitle(osArgs []string) []string {
	forceTitle := regexp.MustCompile(consts.FORCE_TITLE)
	emptyArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && forceTitle.MatchString(anArg) {
			globals.ForceTitle = true
		} else {
			emptyArgs = append(emptyArgs, anArg)
		}
	}
	return emptyArgs
}

func LoadArg(osArgs []string) (string, []string, error) {
	theLoad := consts.DEFAULT_LOAD
	loadPlain := regexp.MustCompile(consts.LOAD_PLAIN)
	loadAndSpeed := regexp.MustCompile(consts.LOAD_AND_SPEED)
	loadChoice := regexp.MustCompile(consts.LOAD_CHOICE)
	loadArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && loadPlain.MatchString(anArg) {
			if loadAndSpeed.MatchString(anArg) {
				theLoad = loadChoice.FindString(anArg)
				if theLoad == "" {
					return "", nil, flaws.BadLoadSerious.MakeFlaw(anArg)
				}
			} else {
				return "", nil, flaws.BadLoadSerious.MakeFlaw(anArg)
			}
		} else {
			loadArgs = append(loadArgs, anArg)
		}
	}
	if theLoad == consts.HIGH_LOAD {
		highMess := "\n\tWARNING, this program will slow down Internet browsing by hogging all bandwith, as it" +
			"\n\tis running under 'high' network load because of the default'--networkLoad=high' option\n"
		fmt.Println(highMess)
	}
	return theLoad, loadArgs, nil
}

func LimitArg(osArgs []string) (int, []string, error) {
	theLimit := consts.DEFAULT_LIMIT
	var err error
	limitPlain := regexp.MustCompile(consts.LIMIT_PLAIN)
	limitAndNumber := regexp.MustCompile(consts.LIMIT_AND_NUMBER)
	limitNumber := regexp.MustCompile(consts.LIMIT_NUMBER)
	limitArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && limitPlain.MatchString(anArg) {
			if limitAndNumber.MatchString(anArg) {
				limitStr := limitNumber.FindString(anArg)
				theLimit, err = strconv.Atoi(limitStr)
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
	return theLimit, limitArgs, nil
}
