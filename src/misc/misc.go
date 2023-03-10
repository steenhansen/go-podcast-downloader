package misc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/models"
)

func EmptyPodcastResults(wasCanceled bool, err error) models.PodcastResults {
	podcastResults := models.PodcastResults{
		ReadFiles:     0,
		SavedFiles:    0,
		PossibleFiles: 0,
		VarietyFiles:  "",
		PodcastTime:   0,
		WasCanceled:   wasCanceled,
		SeriousError:  err,
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
	// logFlag, noLogArgs, err := LoadArg(noLoadArgs)
	// if err != nil {
	// 	panic(err)
	// }
	noEmptyArgs := setEmptyFiles(noLoadArgs)
	noDnsArgs := setDnsErrors(noEmptyArgs)
	noLogArgs := setDnsErrors(noDnsArgs)
	noForceTitle := setForceTitle(noLogArgs)

	cleanArgs := setLogChannels(noForceTitle)
	progPath := CurDir()
	progBounds := models.ProgBounds{
		ProgPath:    progPath,
		LoadOption:  loadFlag,
		LimitOption: limitFlag,
		MinDisk:     minDiskBytes,
	}
	StartLog("/src/" + consts.CHANNEL_LOG_NAME)
	return diskSize, progBounds, cleanArgs
}

func SplitByNewline(multiline string) []string {
	multiline = strings.ReplaceAll(multiline, "\r\n", "\n")
	multiline = strings.ReplaceAll(multiline, "\r", "\n")
	multilines := strings.Split(multiline, "\n")
	return multilines
}

func StartLog(logRelative string) {
	if globals.LogChannels {
		go MemMonitor(consts.MEM_MONITOR_SECONDS)
		progPath := CurDir()
		logPath := progPath + logRelative
		os.Remove(logPath)
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logFile)
		log.Println("--------------------------Start--------------------------")
	}
}

func ChannelLog(channelMess string) {
	if globals.LogChannels {
		log.Println(channelMess)
	}
}

/*

 */
// https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/
// go misc.MemMonitor(300)
func MemMonitor(duration int) {
	var monitorMem models.MonitorMem
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		runtime.ReadMemStats(&rtm)
		monitorMem.Current = rtm.Alloc
		monitorMem.Cumulative = rtm.TotalAlloc
		monitorMem.System = rtm.Sys
		asBytes, _ := json.Marshal(monitorMem)
		ChannelLog("MemMonitor " + string(asBytes))
	}
}
