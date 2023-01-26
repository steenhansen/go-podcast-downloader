package misc

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

func OutputProgress(progressStr string) {
	fmt.Println(progressStr)
}

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

func NameOfFile(furl string) string {
	fparts := strings.Split(furl, "/")
	fname := fparts[len(fparts)-1]
	return fname
}

type VarietiesSet map[string]bool

func (varieties VarietiesSet) AddVariety(filename string) {
	if filename != consts.URL_OF_RSS {
		pieces := strings.Split(filename, ".")
		if len(pieces) > 1 {
			variety := pieces[len(pieces)-1]
			varieties[variety] = true
		}
	}
}

func (varieties VarietiesSet) VarietiesString(separator string) (vString string) {
	for k := range varieties {
		vString = vString + k + " "
	}
	vString = strings.TrimSpace(vString)
	return vString
}

func DiskPanic(fileSize, min_disk_mbs int) error {
	usage := du.NewDiskUsage(".")
	availableUint64 := usage.Available()
	availableBytes := int(availableUint64)
	afterWrite := availableBytes - fileSize
	if afterWrite < min_disk_mbs {
		freeGmb := GbOrMb(afterWrite)
		minimumGmb := GbOrMb(min_disk_mbs)
		freeNeeded := freeGmb + " free need minimum " + minimumGmb + " to proceed"
		err := flaws.LowDisk.StartError(freeNeeded)
		return err
	}
	return nil
}

func DiskSpace() (free, size, percent string) {
	usage := du.NewDiskUsage(".")

	available := usage.Available() / uint64(consts.GB_BYTES)
	free = fmt.Sprintf("%dGB", available)

	capacity := usage.Size() / uint64(consts.GB_BYTES)
	size = fmt.Sprintf("%dGB", capacity)

	used := usage.Usage() * 100
	percent = fmt.Sprintf("%.0f%%", used)
	return free, size, percent
}

func GbOrMbOLD(length int) string {
	if length == 0 {
		return ""
	} else if int64(length) < consts.MB_BYTES {
		lenKb := int64(length) / consts.KB_BYTES
		return fmt.Sprintf("%.0dKB", lenKb)
	} else if int64(length) < consts.GB_BYTES {
		lenMb := int64(length) / consts.MB_BYTES
		return fmt.Sprintf("%.0dMB", lenMb)
	} else {
		lenGb := int64(length) / consts.GB_BYTES
		return fmt.Sprintf("%.0dGB", lenGb)
	}
}

func GbOrMb(length int) string {
	if length == 0 {
		return ""
	} else if int64(length) < consts.MB_BYTES {
		lenKb := int64(length) / consts.KB_BYTES
		return fmt.Sprintf("%.0dKB", lenKb)
	} else if int64(length) < consts.GB_BYTES {
		lenMb := int64(length) / consts.MB_BYTES
		return fmt.Sprintf("%.0dMB", lenMb)
	} else if int64(length) < consts.TB_BYTES {
		lenGb := int64(length) / consts.GB_BYTES
		return fmt.Sprintf("%.0dGB", lenGb)
	} else {
		lenTb := int64(length) / consts.TB_BYTES
		return fmt.Sprintf("%.0dTB", lenTb)
	}
}

func LimitArg(osArgs []string) (int, []string, error) {
	limit := 0
	var err error
	LIMIT_PLAIN := regexp.MustCompile(consts.LIMIT_PLAIN)
	LIMIT_AND_NUMBER := regexp.MustCompile(consts.LIMIT_AND_NUMBER)
	LIMIT_NUMBER := regexp.MustCompile(consts.LIMIT_NUMBER)
	limitArgs := make([]string, 0)
	for i, arg := range osArgs {
		if i > 0 && LIMIT_PLAIN.MatchString(arg) {
			if LIMIT_AND_NUMBER.MatchString(arg) {
				limitStr := LIMIT_NUMBER.FindString(arg)
				limit, err = strconv.Atoi(limitStr)
				if err != nil {
					return 0, nil, err
				}
			} else {
				return 0, nil, flaws.BadLimit.StartError(arg)
			}
		} else {
			limitArgs = append(limitArgs, arg)
		}
	}
	return limit, limitArgs, nil
}
func LoadArg(osArgs []string) (string, []string, error) {
	load := consts.HIGH_LOAD
	LOAD_PLAIN := regexp.MustCompile(consts.LOAD_PLAIN)
	LOAD_AND_SPEED := regexp.MustCompile(consts.LOAD_AND_SPEED)
	LOAD_CHOICE := regexp.MustCompile(consts.LOAD_CHOICE)
	loadArgs := make([]string, 0)
	for i, arg := range osArgs {
		if i > 0 && LOAD_PLAIN.MatchString(arg) {
			if LOAD_AND_SPEED.MatchString(arg) {
				load = LOAD_CHOICE.FindString(arg)
				if load == "" {
					return "", nil, flaws.BadLoad.StartError(arg)
				}
			} else {
				return "", nil, flaws.BadLoad.StartError(arg)
			}
		} else {
			loadArgs = append(loadArgs, arg)
		}
	}
	return load, loadArgs, nil
}

// go run ./ https://sffaudio.herokuapp.com/pdf/rss --limit=3 --load=medium
// go run ./ -race https://sffaudio.herokuapp.com/pdf/rss --limit=3 --load=medium
func DelRace(osArgs []string) ([]string, error) {
	singleDashAlpha := regexp.MustCompile(consts.SINGLE_DASH_ALPHA)
	raceArgs := make([]string, 0)
	for i, arg := range osArgs {
		if i > 0 && singleDashAlpha.MatchString(arg) {
			if arg != consts.RACE_DEBUG {
				return nil, flaws.BadFlag.StartError(arg)
			}
		} else {
			raceArgs = append(raceArgs, arg)
		}
	}
	return raceArgs, nil
}

func IsTesting(osArgs []string) bool {
	for _, arg := range osArgs {
		if strings.HasPrefix(arg, consts.TEST_FLAG_PREFIX) {
			return true
		}
	}
	return false
}
