package media

import (
	"errors"
)

func DoesntExist(containDir, locId string, err error) (string, bool, error) {
	fileStat := errors.New(containDir)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, fileStat, err)
	return "", false, multiErr
}

func CannotCreate(originRss, locId string, err error) (string, bool, error) {
	badHttp := errors.New(originRss)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, badHttp, err)
	return "", false, multiErr
}

func WriteError(mediaUrl, locId string, err error) (string, bool, error) {
	badHttp := errors.New(mediaUrl)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, badHttp, err)
	return "", false, multiErr
}
