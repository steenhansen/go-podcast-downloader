package flaws

import (
	"fmt"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type ErrorKind int

const (
	sStop ErrorKind = iota + 1 // sStop
	badUrl
	timeoutQuit
	emptyRss
	missingTitle
	emptyTitle

	emptyItems
	cantCreateDir
	cantCreateFile
	cantWriteFile
	badChoice
	noPodcasts

	noMatchName
	invalidRssURL
	invalidXML
	firstArgNotUrl

	keyboard
	lowDisk
	badContent
	badLimit
	badLoad
	badFlag
)

var (
	SStop        = flawError{kindError: sStop}
	BadUrl       = flawError{kindError: badUrl}
	TimeoutQuit  = flawError{kindError: timeoutQuit}
	EmptyRss     = flawError{kindError: emptyRss}
	MissingTitle = flawError{kindError: missingTitle}
	EmptyTitle   = flawError{kindError: emptyTitle}

	EmptyItems     = flawError{kindError: emptyItems}
	CantCreateDir  = flawError{kindError: cantCreateDir}
	CantCreateFile = flawError{kindError: cantCreateFile}
	CantWriteFile  = flawError{kindError: cantWriteFile}
	BadChoice      = flawError{kindError: badChoice}
	NoPodcasts     = flawError{kindError: noPodcasts}

	NoMatchName    = flawError{kindError: noMatchName} // NOT USED ??
	InvalidRssURL  = flawError{kindError: invalidRssURL}
	InvalidXML     = flawError{kindError: invalidXML}
	FirstArgNotUrl = flawError{kindError: invalidXML}

	Keyboard = flawError{kindError: keyboard} ///  ??? WHY CONT???? SHOULD BE BASE
	LowDisk  = flawError{kindError: lowDisk}

	BadContent = flawError{kindError: badContent}
	BadLimit   = flawError{kindError: badLimit}
	BadLoad    = flawError{kindError: badLoad}
	BadFlag    = flawError{kindError: badFlag}
)

func (be flawError) Is(otherError error) bool {
	baseKind := be.kindError
	otherKind := otherError.(flawError).kindError
	return baseKind == otherKind
}

func (be flawError) StartError(baseMess string) flawError {
	beNew := be
	beNew.errMess = baseMess
	beNew.previousError = nil
	return beNew
}

func (ce flawError) ContinueError(chainedMess string, startingError error) flawError {
	ceNew := ce
	ceNew.errMess = chainedMess
	ceNew.previousError = startingError
	return ceNew
}

func (ce flawError) Unwrap() error {
	return ce.previousError
}

type flawError struct {
	kindError     ErrorKind
	errMess       string
	previousError error
}

func (e flawError) Error() string {
	pre := consts.ERROR_PREFIX
	switch e.kindError {
	case sStop:
		stopMess := "podcast '%s' was stopped by '" + consts.STOP_KEY_LOWER + "' being entered"
		return fmt.Sprintf(pre+stopMess, e.errMess)
	case badUrl:
		return fmt.Sprintf("XYZA bad url %sCBA", e.errMess) //  go run ./ https://www.naasdfasdfsa.gov
	case timeoutQuit:
		return fmt.Sprintf("Internet timed out %s", e.errMess)
	case emptyRss:
		return fmt.Sprintf("empty rss file %s", e.errMess)

	case missingTitle:
		return "missing title"
	case emptyTitle:
		return "empty title"

	case emptyItems:
		return "empty items"
	case cantCreateDir:
		return fmt.Sprintf("cannot create directory %s", e.errMess)
	case cantCreateFile:
		return fmt.Sprintf("cannot create file %s", e.errMess)
	case cantWriteFile:
		return fmt.Sprintf("cannot write to file %s", e.errMess)

	case badChoice:
		return fmt.Sprintf(pre+"choice does not exist -  %s", e.errMess)
	case noPodcasts:
		return pre + "No podcasts have been added yet, try\n" +
			pre + "$> ./pd-console.exe https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss"

	case noMatchName:
		return fmt.Sprintf("No such podcast name in directory %s", e.errMess)
		// go run ./ www.example.com/
	case invalidRssURL:
		return fmt.Sprintf("Invalid Rss Url %s in %s", e.errMess, consts.URL_OF_RSS)
	case invalidXML:
		return fmt.Sprintf("Invalid XML %s", e.errMess)
	case firstArgNotUrl:
		return fmt.Sprintf("First argument must be the RSS Url, but is instead %s", e.errMess)

	case keyboard:
		return fmt.Sprintf("keyboard error %s", e.errMess)
	case lowDisk:
		return fmt.Sprintf("low disk space, %s", e.errMess)
	case badContent:
		return fmt.Sprintf(pre+"404 html page, %s", e.errMess)
	case badLimit:
		return fmt.Sprintf(pre+"unknown limit option '%s', try '--limit=10' ", e.errMess)
	case badLoad:
		return fmt.Sprintf(pre+"unknown load option '%s', try '--load=low' ", e.errMess)
	case badFlag:
		return fmt.Sprintf(pre+"unknown Go flag '%s', try '-race' ", e.errMess)
	}
	return "unknown error"
}
