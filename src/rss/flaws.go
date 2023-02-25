package rss

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/flaws"
)

func badRetryHttp(retryMess, locId string, err error) (*http.Response, error) {
	locStat := errors.New(locId)
	badHttp := flaws.ExceedRetry.MakeFlaw(retryMess)
	multiErr := errors.Join(locStat, badHttp, err)
	return nil, multiErr
}

func badCallHttp(mediaUrl, locId string, err error) (*http.Response, error) {
	locStat := errors.New(locId)
	badHttp := errors.New(mediaUrl)
	multiErr := errors.Join(locStat, badHttp, err)
	return nil, multiErr
}

func badHttp(mediaUrl, locId string, err error) (*http.Response, error) {
	locStat := errors.New(locId)
	badHttp := errors.New(mediaUrl)
	multiErr := errors.Join(locStat, badHttp, err)
	return nil, multiErr
}

func not200Flaw(respStat, mediaUrl, locId string) (int, error) {
	locStat := errors.New(locId)
	err := flaws.HttpFault.MakeFlaw(respStat + consts.ERROR_SEPARATOR + mediaUrl)
	multiErr := errors.Join(locStat, err)
	return 0, multiErr
}

func readAllFlaw(mediaUrl, locId string, err error) (int, error) {
	locStat := errors.New(locId)
	badUrl := errors.New(mediaUrl)
	multiErr := errors.Join(locStat, err, badUrl)
	return 0, multiErr
}

func osCreateFlaw(filePath, locId string, err error) (int, error) {
	locStat := errors.New(locId)
	fileStat := errors.New(filePath)
	multiErr := errors.Join(locStat, fileStat, err)
	return 0, multiErr
}
func was404Flaw(filePath, mediaUrl, locId string, err error) (int, error) {

	os.Remove(filePath)
	locStat := errors.New(locId)
	htmlStat := errors.New("not a media file, instead html starting with " + consts.HTML_404_BEGIN)
	badUrl := errors.New(mediaUrl)
	multiErr := errors.Join(locStat, htmlStat, badUrl, err)
	return 0, multiErr
}

func badWriteFlaw(mediaFile *os.File, filePath, locId string, err error) (int, error) {
	mediaFile.Close()
	os.Remove(filePath)
	locStat := errors.New(locId)
	fileStat := errors.New(filePath)
	multiErr := errors.Join(locStat, fileStat, err)
	return 0, multiErr
}

func diskPanicFlaw(mediaFile *os.File, filePath, locId string, err error) (int, error) {
	mediaFile.Close()
	os.Remove(filePath)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, err)
	return 0, multiErr
}

func length0Flaw(mediaFile *os.File, filePath, locId string) (int, error) {
	mediaFile.Close()
	os.Remove(filePath)
	fileStat := flaws.EmptyFileWrite.MakeFlaw(filePath)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, fileStat)
	return 0, multiErr
}
func lengthWrongFlaw(mediaFile *os.File, filePath, locId string, writtenBytes, lengthContent int) (int, error) {
	mediaFile.Close()
	os.Remove(filePath)
	writtenStr := strconv.Itoa(writtenBytes)
	lengthStr := strconv.Itoa(lengthContent)
	fileWrittenLength := filePath + ", " + writtenStr + " bytes, " + lengthStr + " bytes"
	fileStat := flaws.InvalidFileWrite.MakeFlaw(fileWrittenLength)
	locStat := errors.New(locId)
	multiErr := errors.Join(locStat, fileStat)
	return 0, multiErr
}
