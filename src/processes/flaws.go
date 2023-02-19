package processes

import (
	"context"
	"errors"
	"fmt"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
)

func keyboardFailure(err error) models.PodcastResults {
	locStat := errors.New("DOWNLOADMEDIA()")
	multiErr := errors.Join(locStat, err)
	emptyPodcastResults := misc.EmptyPodcastResults(multiErr)
	return emptyPodcastResults

}

func firstErr(err error, seriousStream <-chan error) error {
	for i := 0; i < len(seriousStream); i++ {
		err := <-seriousStream
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func dealWithErrors(err, ctxErr error, podcastData models.PodcastData, podcastResults models.PodcastResults) models.PodcastResults {
	if globals.StopingOnSKey {
		podcastResults.Err = flaws.SKeyStop.MakeFlaw(podcastData.PodTitle)
		return podcastResults
	}

	//	fmt.Println("dealWithErrors looking for timeoutStop/lowDisk")
	//	fmt.Println("err=", err)
	// fmt.Println("ctxErr=", ctxErr)
	if err != nil {
		return misc.EmptyPodcastResults(err)
	}

	// // so was context canceled for a DNS error, a diskSpace error, or an 'S'
	//fmt.Println("dealWithErrors looking for timeoutStop/lowDisk ctxErr=", ctxErr)
	if ctxErr == context.Canceled {
		fmt.Println("context was canceled")
		return podcastResults
	}

	if ctxErr != nil {
		fmt.Println("context error, not canceled, but instead ctxErr=", ctxErr)
		return misc.EmptyPodcastResults(ctxErr)
	}
	return podcastResults
}
