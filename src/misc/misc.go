package misc

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func EmptyPodcastResults(err error) consts.PodcastResults {
	podcastResults := consts.PodcastResults{
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
		loadProcs = maxProcessors / 2
	}
	return loadProcs
}

func DiskPanic(fileSize, minDiskMbs int) error {
	dUsage := du.NewDiskUsage(".")
	availableUint64 := dUsage.Available()
	availableBytes := int(availableUint64)
	afterWrite := availableBytes - fileSize
	if afterWrite < minDiskMbs {
		freeGmb := GbOrMb(afterWrite)
		minimumGmb := GbOrMb(minDiskMbs)
		freeNeeded := freeGmb + " free need minimum " + minimumGmb + " to proceed"
		err := flaws.LowDiskSerious.StartError(freeNeeded)
		return err
	}
	return nil
}

func GbOrMb(dirSize int) string {
	if dirSize == 0 {
		return ""
	} else if int64(dirSize) < consts.MB_BYTES {
		lenKb := int64(dirSize) / consts.KB_BYTES
		return fmt.Sprintf("%.0dKB", lenKb)
	} else if int64(dirSize) < consts.GB_BYTES {
		lenMb := int64(dirSize) / consts.MB_BYTES
		return fmt.Sprintf("%.0dMB", lenMb)
	} else if int64(dirSize) < consts.TB_BYTES {
		lenGb := int64(dirSize) / consts.GB_BYTES
		return fmt.Sprintf("%.0dGB", lenGb)
	} else {
		lenTb := int64(dirSize) / consts.TB_BYTES
		return fmt.Sprintf("%.0dTB", lenTb)
	}
}

func LimitArg(osArgs []string) (int, []string, error) {
	theLimit := 0
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
				return 0, nil, flaws.BadLimitSerious.StartError(anArg)
			}
		} else {
			limitArgs = append(limitArgs, anArg)
		}
	}
	return theLimit, limitArgs, nil
}

func LoadArg(osArgs []string) (string, []string, error) {
	theLoad := consts.HIGH_LOAD
	loadPlain := regexp.MustCompile(consts.LOAD_PLAIN)
	loadAndSpeed := regexp.MustCompile(consts.LOAD_AND_SPEED)
	loadChoice := regexp.MustCompile(consts.LOAD_CHOICE)
	loadArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && loadPlain.MatchString(anArg) {
			if loadAndSpeed.MatchString(anArg) {
				theLoad = loadChoice.FindString(anArg)
				if theLoad == "" {
					return "", nil, flaws.BadLoadSerious.StartError(anArg)
				}
			} else {
				return "", nil, flaws.BadLoadSerious.StartError(anArg)
			}
		} else {
			loadArgs = append(loadArgs, anArg)
		}
	}
	return theLoad, loadArgs, nil
}

// go run ./ https://sffaudio.herokuapp.com/pdf/rss --limit=3 --load=medium
// go run ./ -race https://sffaudio.herokuapp.com/pdf/rss --limit=3 --load=medium
func DelRace(osArgs []string) ([]string, error) {
	singleDashAlpha := regexp.MustCompile(consts.SINGLE_DASH_ALPHA)
	raceArgs := make([]string, 0)
	for argIndex, anArg := range osArgs {
		if argIndex > 0 && singleDashAlpha.MatchString(anArg) {
			if anArg != consts.RACE_DEBUG {
				return nil, flaws.BadFlagSerious.StartError(anArg)
			}
		} else {
			raceArgs = append(raceArgs, anArg)
		}
	}
	return raceArgs, nil
}

func KeyboardMenuChoice() string {
	keyboardReader := bufio.NewReader(os.Stdin)
	inputText, _ := keyboardReader.ReadString('\n')
	return inputText
}

func IsTesting(osArgs []string) bool {
	for _, anArg := range osArgs {
		if strings.HasPrefix(anArg, consts.TEST_FLAG_PREFIX) {
			return true
		}
	}
	return false
}

func CurDir() string {
	progPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return progPath
}

type diskSpaceFn func() (string, string, string)

func DiskSpace() (dFree, dSize, dPercent string) {
	dUsage := du.NewDiskUsage(".")

	dAvailable := dUsage.Available() / uint64(consts.GB_BYTES)
	dFree = fmt.Sprintf("%dGB", dAvailable)

	dCapacity := dUsage.Size() / uint64(consts.GB_BYTES)
	dSize = fmt.Sprintf("%dGB", dCapacity)

	dUsed := dUsage.Usage() * 100
	dPercent = fmt.Sprintf("%.0f%%", dUsed)
	return dFree, dSize, dPercent
}

func InitProg(diskSpace diskSpaceFn, minDiskBytes int) (string, consts.ProgBounds, []string) {
	dFree, dSize, dPercent := diskSpace()
	diskSize := fmt.Sprintf("Current disk has %s free from a total %s which is %s full\n", dFree, dSize, dPercent)
	raceArgs, err := DelRace(os.Args)
	if err != nil {
		panic(err)
	}
	limitFlag, tempArgs, err := LimitArg(raceArgs)
	if err != nil {
		panic(err)
	}
	loadFlag, cleanArgs, err := LoadArg(tempArgs)
	if err != nil {
		panic(err)
	}

	progPath := CurDir()
	progBounds := consts.ProgBounds{
		ProgPath:    progPath,
		LoadOption:  loadFlag,
		LimitOption: limitFlag,
		MinDisk:     minDiskBytes,
	}

	return diskSize, progBounds, cleanArgs
}
