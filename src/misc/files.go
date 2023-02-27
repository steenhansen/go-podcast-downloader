package misc

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/ricochet2200/go-disk-usage/du"
	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
)

func FileTimeout(maxReadFileTime time.Duration) time.Duration {
	return maxReadFileTime * consts.HTTP_RETRIES
}

func FilesInDir(dirPath string) ([]fs.FileInfo, error) {
	podDir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer podDir.Close()
	dirFiles, err := podDir.Readdir(0)
	return dirFiles, err
}

func DiskPanic(fileSize, minDiskMbs int) error {
	dUsage := du.NewDiskUsage(".")
	availableUint64 := dUsage.Available()
	availableBytes := int(availableUint64)
	afterWrite := availableBytes - fileSize
	if afterWrite < minDiskMbs {
		freeGmb := GbOrMb(afterWrite)
		minimumGmb := GbOrMb(minDiskMbs)
		freeNeeded := freeGmb + " free, need minimum " + minimumGmb + " to proceed"
		err := flaws.LowDiskSerious.MakeFlaw(freeNeeded)
		return err
	}
	return nil
}

func GbOrMb(dirSize int) string {
	if int64(dirSize) < consts.KB_BYTES {
		lenB := int64(dirSize)
		if lenB == 0 {
			return ""
		}
		return fmt.Sprintf("%.0dB", lenB)
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

func CurDir() string {
	progPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return progPath
}

func diskSpace() (dFree, dSize, dPercent string) {
	dUsage := du.NewDiskUsage(".")

	dAvailable := dUsage.Available() / uint64(consts.GB_BYTES)
	dFree = fmt.Sprintf("%dGB", dAvailable)

	dCapacity := dUsage.Size() / uint64(consts.GB_BYTES)
	dSize = fmt.Sprintf("%dGB", dCapacity)

	dUsed := dUsage.Usage() * 100
	dPercent = fmt.Sprintf("%.0f%%", dUsed)
	return dFree, dSize, dPercent
}
