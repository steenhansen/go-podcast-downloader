package processes

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
)

func GoDownloadError(ctx context.Context, cancel context.CancelFunc, errorStream <-chan error, toFixStream chan<- error) {
osErrCancel:
	for {
		select {
		case <-ctx.Done():
			break osErrCancel
		case err := <-errorStream:
			if !errors.Is(err, flaws.BadUrl) && !errors.Is(err, flaws.BadContent) {
				toFixStream <- err
				cancel()
				break osErrCancel ///errors.Is(err, flaws.LowDisk)
			}
		}
	}
}

// https://github.com/eiannone/keyboard
func GoStopKey(ctx context.Context, cancel context.CancelFunc, dirName string, keysEvents <-chan keyboard.KeyEvent, simKeyStream <-chan string) {
	misc.OutputProgress("Downloading '" + dirName + "' podcast, hit '" + consts.STOP_KEY_LOWER + "' to stop")
sKeyCancel:
	for {
		select {
		case event := <-keysEvents:
			keyChar := string(event.Rune)
			keyLower := strings.ToLower(keyChar)
			if keyLower == consts.STOP_KEY_LOWER {
				cancel()
				break sKeyCancel
			}
		case simKey := <-simKeyStream:
			misc.OutputProgress("TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			cancel()
			break sKeyCancel
		case <-ctx.Done():
			break sKeyCancel
		}
	}
}

// readFiles, savedFiles *int, min_disk_mbs int,
func GoDownloadAndSaveFiles(ctx context.Context, mediaStream <-chan consts.UrlPathLength, dtoFixStream chan<- bool, readFiles, savedFiles *int, min_disk_mbs int,
	errorStream chan<- error) {
downloadCancel:
	for media := range mediaStream {
		select {
		case <-ctx.Done():
			break downloadCancel
		default:
			start := time.Now()
			misc.OutputProgress(feed.ShowProgress(media, readFiles))

			// misc.OutputProgress() needs to save to a buffer

			err := feed.DownloadAndWriteFile(ctx, media.Url, media.Path, min_disk_mbs)
			if err != nil {
				errorStream <- err
				misc.MediaFaults2(media.Url, err)
				misc.OutputProgress(feed.ShowError(media.Url))
			} else if ctx.Err() == nil {
				misc.OutputProgress(feed.ShowSaved(savedFiles, start, media.Url))
			}
		}
	}
	dtoFixStream <- true
}

func DownloadMedia(url string, PodcastData consts.PodcastData, progBounds consts.ProgBounds, simKeyStream chan string) consts.PodcastResults {
	misc.OutputProgress(consts.CLEAR_SCREEN)
	allTime := time.Now()
	workers := misc.NumWorkers(progBounds.LoadOption)
	mediaStream := make(chan consts.UrlPathLength)
	dtoFixStream := make(chan bool, workers)
	errorStream := make(chan error, workers)
	toFixStream := make(chan error, workers)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(consts.MAX_READ_FILE_TIME))
	defer cancel()
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return misc.EmptyPodcastResults(flaws.Keyboard.ContinueError("cc", err))
	}
	var readFiles, savedFiles int
	go GoDownloadError(ctx, cancel, errorStream, toFixStream)
	go GoStopKey(ctx, cancel, PodcastData.MediaTitle, keysEvents, simKeyStream)
	for i := 0; i < workers; i++ {
		go GoDownloadAndSaveFiles(ctx, mediaStream, dtoFixStream, &readFiles, &savedFiles, progBounds.MinDisk, errorStream)
	}
	possibleFiles, varietyFiles, err := media.SaveDownloadedMedia(ctx, PodcastData, mediaStream, progBounds.LimitOption)
	close(mediaStream)
	for i := 0; i < workers; i++ {
		<-dtoFixStream
	}
	err = firstErr(err, toFixStream)
	if err != nil {
		return misc.EmptyPodcastResults(err)
	}
	podcastResults := consts.PodcastResults{
		ReadFiles:     readFiles,
		SavedFiles:    savedFiles,
		PossibleFiles: possibleFiles,
		VarietyFiles:  varietyFiles,
		PodcastTime:   time.Since(allTime),
		Err:           nil,
	}
	if ctx.Err() == context.Canceled {
		podcastResults.Err = flaws.SStop.StartError(PodcastData.MediaPath)
		return podcastResults
		//return misc.EmptyPodcastResults(flaws.SStop.StartError(PodcastData.MediaPath))
	}
	if ctx.Err() != nil {
		return misc.EmptyPodcastResults(ctx.Err())
	}

	return podcastResults
}

func firstErr(err error, toFixStream <-chan error) error {
	if err != nil {
		return err
	}
	for i := 0; i < len(toFixStream); i++ {
		err := <-toFixStream
		if err != nil {
			return err
		}
	}
	return nil
}
