package flaws

import (
	"errors"
	"fmt"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

type ErrorKind int

const (
	sStop ErrorKind = iota + 1

	badFlagSerious
	badLimitSerious
	badLoadSerious
	cantCreateDirSerious
	cantCreateFileSerious
	cantWriteFileSerious
	keyboardSerious
	lowDiskSerious

	badChoice
	badContent
	badUrl
	emptyItems
	emptyRss
	emptyTitle
	firstArgNotUrl
	invalidRssURL
	invalidXML
	missingTitle
	noMatchName
	noPodcasts
	timeoutQuit
)

type flawError struct {
	kindError     ErrorKind
	errMess       string
	previousError error
}

func IsSerious(err error) bool {
	if err != nil {
		is_serious := errors.Is(err, LowDiskSerious) ||
			errors.Is(err, CantCreateDirSerious) ||
			errors.Is(err, CantCreateFileSerious) ||
			errors.Is(err, CantWriteFileSerious)
		return is_serious
	}
	return false
}

var (
	SStop = flawError{kindError: sStop}

	BadFlagSerious        = flawError{kindError: badFlagSerious}
	BadLimitSerious       = flawError{kindError: badLimitSerious}
	BadLoadSerious        = flawError{kindError: badLoadSerious}
	CantCreateDirSerious  = flawError{kindError: cantCreateDirSerious}
	CantCreateFileSerious = flawError{kindError: cantCreateFileSerious}
	CantWriteFileSerious  = flawError{kindError: cantWriteFileSerious}
	KeyboardSerious       = flawError{kindError: keyboardSerious}
	LowDiskSerious        = flawError{kindError: lowDiskSerious}

	BadChoice      = flawError{kindError: badChoice}
	BadContent     = flawError{kindError: badContent}
	BadUrl         = flawError{kindError: badUrl}
	EmptyItems     = flawError{kindError: emptyItems}
	EmptyRss       = flawError{kindError: emptyRss}
	EmptyTitle     = flawError{kindError: emptyTitle}
	FirstArgNotUrl = flawError{kindError: invalidXML}
	InvalidRssURL  = flawError{kindError: invalidRssURL}
	InvalidXML     = flawError{kindError: invalidXML}
	MissingTitle   = flawError{kindError: missingTitle}
	NoMatchName    = flawError{kindError: noMatchName} // in go run ./ my-nasa   is not there!
	NoPodcasts     = flawError{kindError: noPodcasts}
	TimeoutQuit    = flawError{kindError: timeoutQuit}
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

func (e flawError) Error() string {
	pre := consts.ERROR_PREFIX
	switch e.kindError {
	case sStop:
		stopMess := "podcast '%s' was stopped by '" + consts.STOP_KEY_LOWER + "' being entered"
		return fmt.Sprintf(pre+stopMess, e.errMess)

	case badFlagSerious:
		return fmt.Sprintf(pre+"unknown Go flag '%s', try '-race' ", e.errMess)
	case badLimitSerious:
		return fmt.Sprintf(pre+"unknown limit option '%s', try '--limit=10' ", e.errMess)
	case badLoadSerious:
		return fmt.Sprintf(pre+"unknown load option '%s', try '--load=low' ", e.errMess)
	case cantCreateDirSerious:
		return fmt.Sprintf(pre+"cannot create directory %s", e.errMess)
	case cantCreateFileSerious:
		return fmt.Sprintf(pre+"cannot create file %s", e.errMess)
	case cantWriteFileSerious:
		return fmt.Sprintf(pre+"cannot write to file %s", e.errMess)
	case keyboardSerious:
		return fmt.Sprintf(pre+"keyboard error %s", e.errMess)
	case lowDiskSerious:
		return fmt.Sprintf(pre+"low disk space, %s", e.errMess)

	case badChoice:
		return fmt.Sprintf(pre+"choice does not exist -  %s", e.errMess)
	case badContent:
		return fmt.Sprintf(pre+"404 or 400 html page, %s", e.errMess)
	case badUrl:
		return fmt.Sprintf(pre+"bad url %s", e.errMess) //  go run ./ https://www.naasdfasdfsa.gov
	case emptyItems:
		return pre + "empty items"
	case emptyRss:
		return fmt.Sprintf(pre+"empty rss file %s", e.errMess)
	case emptyTitle:
		return pre + "empty title"
	case firstArgNotUrl:
		return fmt.Sprintf(pre+"First argument must be the RSS Url, but is instead %s", e.errMess)
	case invalidRssURL:
		return fmt.Sprintf(pre+"Invalid Rss Url %s in %s", e.errMess, consts.URL_OF_RSS_FN)
	case invalidXML:
		return fmt.Sprintf(pre+"Invalid XML %s", e.errMess)
	case missingTitle:
		return pre + "missing title"
	case noMatchName:
		return fmt.Sprintf(pre+"The podcast folder '%s' was not found", e.errMess) // go run ./ my-nasa
	case noPodcasts:
		return pre + "No podcasts have been added yet, try\n" +
			pre + "$> ./pd-console.exe https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss"
	case timeoutQuit:
		return fmt.Sprintf(pre+"Internet timed out %s", e.errMess)
	}
	return pre + "unknown error"
}
