package processes

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/models"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func GoDownloadError(ctx context.Context, cancel context.CancelFunc, errorStream <-chan error, seriousStream chan<- error) {
osErrCancel:
	for {
		select {
		case <-ctx.Done():
			break osErrCancel
		case err := <-errorStream:
			if flaws.IsSerious(err) { // don't crash on a missing media file
				seriousStream <- err
				cancel()
				break osErrCancel
			}
		}
	}
}

// keyStream <-chan string is so that a test can simulate stopping
func startDownloading(cancel context.CancelFunc, podcastData models.PodcastData) {
	theTitle := podcastData.PodTitle
	numFiles := strconv.Itoa(len(podcastData.PodUrls))
	stopKey := consts.STOP_KEY_LOWER
	globals.Console.Note("Downloading '" + theTitle + "' podcast, " + numFiles + " files, hit '" + stopKey + "' to stop\n")
	if globals.EmptyFilesTest {
		dirFiles, err := misc.FilesInDir(podcastData.PodPath)
		if err != nil {
			cancel()
		}
		if len(dirFiles) > 1 {
			fmt.Println("WARNING '" + podcastData.PodPath + "' should have no media files to test correctly")
		}
	}
}

func GoDownloadAndSaveFiles(ctx context.Context, curStat models.CurStat, mediaStream <-chan models.MediaEnclosure, doneStream chan<- bool, errorStream chan<- error, httpMedia models.HttpFn) {
downloadCancel:
	for media := range mediaStream {
		misc.SleepTime(curStat.NetworkLoad)
		select {
		case <-ctx.Done():
			break downloadCancel
		default:
			start := time.Now()
			globals.Console.Note(feed.ShowProgress(media, curStat.ReadFiles))
			writenBytes, err := rss.DownloadAndWriteFile(ctx, media.EnclosureUrl, media.EnclosurePath, curStat.MinDiskMbs, httpMedia)
			if err != nil {
				if ctx.Err() != context.Canceled {
					errorStream <- err
					globals.Faults.Note(media.EnclosurePath, err)
					globals.Console.Note(feed.ShowError(media.EnclosurePath))
				}
			} else {
				globals.Console.Note(feed.ShowSaved(curStat.SavedFiles, start, media.EnclosurePath))
				if media.EnclosureSize != writenBytes {
					globals.Console.Note(feed.ShowSizeError(media.EnclosureSize, writenBytes))
				}
			}
		}
	}
	doneStream <- true
}

func disposeDownloaders(numWorkers int, doneStream <-chan bool) {
	for i := 0; i < numWorkers; i++ {
		<-doneStream
	}
}

func DownloadMedia(url string, podcastData models.PodcastData, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) models.PodcastResults {
	startTime := time.Now()
	globals.StopingOnSKey = false
	numWorkers := misc.NumWorkers(progBounds.LoadOption)
	mediaStream := make(chan models.MediaEnclosure)
	doneStream := make(chan bool, numWorkers)
	errorStream := make(chan error, numWorkers)
	seriousStream := make(chan error, numWorkers)
	timeOut := misc.FileTimeout(consts.MEDIA_MAX_READ_FILE_TIME)
	ctxMedias, cancelMedias := context.WithTimeout(context.Background(), timeOut)
	defer cancelMedias()
	KeyEventsReal, err := keyboard.GetKeys(consts.KEY_BUFF_SIZE)
	if err != nil {
		return keyboardFailure(err)
	}
	var readFiles, savedFiles int
	go GoDownloadError(ctxMedias, cancelMedias, errorStream, seriousStream)
	startDownloading(cancelMedias, podcastData)
	go misc.GoStopKey(ctxMedias, cancelMedias, KeyEventsReal, keyStream)
	curStat := models.CurStat{
		ReadFiles:   &readFiles,
		SavedFiles:  &savedFiles,
		MinDiskMbs:  progBounds.MinDisk,
		NetworkLoad: progBounds.LoadOption,
	}
	for i := 0; i < numWorkers; i++ {
		go GoDownloadAndSaveFiles(ctxMedias, curStat, mediaStream, doneStream, errorStream, httpMedia)
	}
	possibleFiles, err := media.SaveDownloadedMedia(ctxMedias, podcastData, mediaStream, progBounds.LimitOption, httpMedia) // fix orders like above !!!
	if err != nil {
		cancelMedias()
	}
	close(mediaStream)
	disposeDownloaders(numWorkers, doneStream)
	err = firstErr(err, seriousStream)
	podcastResults := models.PodcastResults{
		ReadFiles:     readFiles,
		SavedFiles:    savedFiles,
		PossibleFiles: possibleFiles,
		PodcastTime:   time.Since(startTime),
		Err:           nil,
	}
	resultsWithErr := dealWithErrors(err, ctxMedias.Err(), podcastData, podcastResults)
	return resultsWithErr
}
