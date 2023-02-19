package media

import (
	"errors"
)

func doesntExist(containDir, locId string, err error) (string, bool, error) {
	fileStat := errors.New(containDir)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, fileStat, err)
	return "", false, multiErr
}

func cannotCreate(originRss, locId string, err error) (string, bool, error) {
	badHttp := errors.New(originRss)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, badHttp, err)
	return "", false, multiErr
}

func writeError(mediaUrl, locId string, err error) (string, bool, error) {
	badHttp := errors.New(mediaUrl)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, badHttp, err)
	return "", false, multiErr
}
