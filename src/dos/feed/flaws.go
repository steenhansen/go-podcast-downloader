package feed

import "errors"

func badStatusCode(mediaUrl, statusCode, locId string, err error) ([]byte, error) {
	mediaErr := errors.New(mediaUrl)
	statusErr := errors.New(statusCode)
	locStat := errors.New(locId)
	multiErr := errors.Join(mediaErr, statusErr, locStat, err)
	return nil, multiErr
}

func badReadAll(mediaUrl, locId string, err error) ([]byte, error) {
	urlErr := errors.New(mediaUrl)
	locStat := errors.New(locId)
	multiErr := errors.Join(urlErr, locStat, err)
	return nil, multiErr
}

func readZero(mediaUrl, locId string, err error) ([]byte, error) {
	urlErr := errors.New(mediaUrl)
	zeroErr := errors.New("read zero bytes")
	locStat := errors.New(locId)
	multiErr := errors.Join(urlErr, zeroErr, locStat, err)
	return nil, multiErr
}
