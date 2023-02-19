package misc

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
)

func EmptyPodcastResults(err error) models.PodcastResults {
	podcastResults := models.PodcastResults{
		ReadFiles:     0,
		SavedFiles:    0,
		PossibleFiles: 0,
		VarietyFiles:  "",
		PodcastTime:   0,
		Err:           err,
	}
	return podcastResults
}

func NumWorkers(loadFlag string) int {
	maxProcessors := runtime.GOMAXPROCS(0)
	loadProcs := 1
	if loadFlag == consts.HIGH_LOAD {
		loadProcs = maxProcessors
	} else if loadFlag == consts.MEDIUM_LOAD {
		loadProcs = maxProcessors / consts.CPUS_MED_DIVIDER
	}
	return loadProcs
}

func SleepTime(loadFlag string) {
	sleepTime := consts.LOW_SLEEP
	if loadFlag == consts.HIGH_LOAD {
		sleepTime = consts.HIGH_SLEEP
	} else if loadFlag == consts.MEDIUM_LOAD {
		sleepTime = consts.MEDIUM_SLEEP
	}
	sleepDuration := time.Duration(sleepTime) * time.Second
	time.Sleep(sleepDuration)
}

func InitProg(minDiskBytes int) (string, models.ProgBounds, []string) {
	dFree, dSize, dPercent := diskSpace()
	diskSize := fmt.Sprintf("Current disk has %s free from a total %s which is %s full\n", dFree, dSize, dPercent)
	raceArgs, err := DelRace(os.Args)
	if err != nil {
		panic(err)
	}
	limitFlag, noLimitArgs, err := LimitArg(raceArgs)
	if err != nil {
		panic(err)
	}
	loadFlag, noLoadArgs, err := LoadArg(noLimitArgs)
	if err != nil {
		panic(err)
	}
	noEmptyArgs := setEmptyFiles(noLoadArgs)
	noDnsArgs := setDnsErrors(noEmptyArgs)
	cleanArgs := setForceTitle(noDnsArgs)
	progPath := CurDir()
	progBounds := models.ProgBounds{
		ProgPath:    progPath,
		LoadOption:  loadFlag,
		LimitOption: limitFlag,
		MinDisk:     minDiskBytes,
	}
	return diskSize, progBounds, cleanArgs
}

func SplitByNewline(multiline string) []string {
	multiline = strings.ReplaceAll(multiline, "\r\n", "\n")
	multiline = strings.ReplaceAll(multiline, "\r", "\n")
	multilines := strings.Split(multiline, "\n")
	return multilines
}
