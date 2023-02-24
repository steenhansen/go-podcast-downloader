package flaws

import (
	"errors"
	"fmt"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
)

const (
	sKeyStop errorKind = iota + 1
	timeoutStop
	httpFault
	exceedRetry

	badChoice
	badFlagSerious
	badLimitSerious
	badLoadSerious
	lowDiskSerious

	emptyItems
	emptyTitle
	emptyPodcasts
	emptyFileWrite

	invalidRssURL
	invalidXML
	invalidXmlTitle
	invalidPodcastName
	invalidFileWrite
)

type errorKind int

type FlawError struct {
	errMess   string
	kindError errorKind
	errs      []error
}

func (flaw FlawError) Unwrap() []error {
	return flaw.errs
}

func IsSerious(err error) bool {
	if err != nil {
		is_serious := errors.Is(err, LowDiskSerious)
		return is_serious
	}
	return false
}

var (
	SKeyStop    = FlawError{kindError: sKeyStop}
	TimeoutStop = FlawError{kindError: timeoutStop}
	HttpFault   = FlawError{kindError: httpFault}
	ExceedRetry = FlawError{kindError: exceedRetry}

	BadChoice       = FlawError{kindError: badChoice}
	BadFlagSerious  = FlawError{kindError: badFlagSerious}
	BadLimitSerious = FlawError{kindError: badLimitSerious}
	BadLoadSerious  = FlawError{kindError: badLoadSerious}
	LowDiskSerious  = FlawError{kindError: lowDiskSerious}

	EmptyItems     = FlawError{kindError: emptyItems}
	EmptyTitle     = FlawError{kindError: emptyTitle}
	EmptyPodcasts  = FlawError{kindError: emptyPodcasts}
	EmptyFileWrite = FlawError{kindError: emptyFileWrite}

	InvalidRssURL      = FlawError{kindError: invalidRssURL}
	InvalidXML         = FlawError{kindError: invalidXML}
	InvalidXmlTitle    = FlawError{kindError: invalidXmlTitle}
	InvalidPodcastName = FlawError{kindError: invalidPodcastName}

	InvalidFileWrite = FlawError{kindError: invalidFileWrite}
)

func (flaw FlawError) Is(otherError error) bool {
	if otherAsFlaw, ok := otherError.(FlawError); ok {
		baseKind := flaw.kindError
		otherKind := otherAsFlaw.kindError
		return baseKind == otherKind
	}
	return false
}

func (flaw FlawError) MakeFlaw(baseMess string) FlawError {
	newFlaw := flaw
	newFlaw.errMess = baseMess
	return newFlaw
}

func (flaw FlawError) Error() string {
	switch flaw.kindError {
	case sKeyStop:
		stopMess := "podcast '%s' was stopped by the '" + consts.STOP_KEY_LOWER + "' key being pressed"
		return fmt.Sprintf(stopMess, flaw.errMess)
	case timeoutStop:
		return fmt.Sprintf("Internet timed out %s", flaw.errMess)
	case httpFault:
		return fmt.Sprint("HTTP error " + flaw.errMess)
	case exceedRetry:
		return fmt.Sprintf("exceeded allowed retries : %s", flaw.errMess)

	case badChoice:
		return fmt.Sprintf("choice does not exist -  %s", flaw.errMess)
	case badFlagSerious:
		return fmt.Sprintf("unknown Go flag '%s', try '-race' ", flaw.errMess)
	case badLimitSerious:
		return fmt.Sprintf("unknown limit option '%s', try '--limit=10' ", flaw.errMess)
	case badLoadSerious:
		return fmt.Sprintf("unknown load option '%s', try '--load=low' ", flaw.errMess)
	case lowDiskSerious:
		return fmt.Sprintf("low disk space, %s", flaw.errMess)

	case emptyItems:
		return "empty items"
	case emptyTitle:
		return "empty title"
	case emptyPodcasts:
		return "No podcasts have been added yet, try\n" +
			"$> ./pd-consolflaw.exe https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss"
	case emptyFileWrite:
		return fmt.Sprintf("empty written file : %s", flaw.errMess)

	case invalidRssURL:
		return fmt.Sprintf("Invalid Rss Url %s in %s", flaw.errMess, consts.URL_OF_RSS_FN)
	case invalidXML:
		return fmt.Sprintf("Invalid XML %s", flaw.errMess)
	case invalidXmlTitle:
		return "missing title"
	case invalidPodcastName:
		return fmt.Sprintf("The podcast folder '%s' was not found", flaw.errMess)

	case invalidFileWrite:
		return fmt.Sprintf("wrong number of bytes written : %s", flaw.errMess)

	}
	return "unknown error"
}
