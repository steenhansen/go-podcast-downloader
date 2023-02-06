package processes

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/globals"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func GoDownloadError(ctx context.Context, cancel context.CancelFunc, errorStream <-chan error, seriousStream chan<- error) {
osErrCancel:
	for {
		select {
		case <-ctx.Done():
			break osErrCancel
		case err := <-errorStream:
			if flaws.IsSerious(err) {
				seriousStream <- err
				cancel()
				break osErrCancel
			}
		}
	}
}

// keyStream <-chan string is so that a test can simulate stopping
func GoStopKey(ctx context.Context, cancel context.CancelFunc, podcastData consts.PodcastData, keysEvents <-chan keyboard.KeyEvent, keyStream <-chan string) {
	theTitle := podcastData.PodTitle
	numFiles := strconv.Itoa(len(podcastData.PodUrls))
	stopKey := consts.STOP_KEY_LOWER
	globals.Console.Note("Downloading '" + theTitle + "' podcast, " + numFiles + " files, hit '" + stopKey + "' to stop\n")
keyboardCancel:
	for {
		select {
		case event := <-keysEvents:
			keyChar := string(event.Rune)
			keyLower := strings.ToLower(keyChar)
			if keyLower == consts.STOP_KEY_LOWER {
				cancel()
				break keyboardCancel
			}
		case simKey := <-keyStream:
			globals.Console.Note("TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			cancel()
			break keyboardCancel
		case <-ctx.Done():
			break keyboardCancel
		}
	}
}

func GoDownloadAndSaveFiles(ctx context.Context, curStat consts.CurStat, mediaStream <-chan consts.MediaEnclosure, doneStream chan<- bool, errorStream chan<- error, httpMedia consts.HttpFunc) {
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
					globals.Faults.Note(media.EnclosureUrl, err)
					globals.Console.Note(feed.ShowError(media.EnclosureUrl))
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

func DownloadMedia(url string, podcastData consts.PodcastData, progBounds consts.ProgBounds, keyStream chan string, httpMedia consts.HttpFunc) consts.PodcastResults {
	startTime := time.Now()
	numWorkers := misc.NumWorkers(progBounds.LoadOption)
	mediaStream := make(chan consts.MediaEnclosure)
	doneStream := make(chan bool, numWorkers)
	errorStream := make(chan error, numWorkers)
	seriousStream := make(chan error, numWorkers)
	ctxMedias, cancelMedias := context.WithTimeout(context.Background(), time.Duration(consts.MAX_READ_FILE_TIME))
	defer cancelMedias()
	keysEvents, err := keyboard.GetKeys(consts.KEY_BUFF_SIZE)
	if err != nil {
		return misc.EmptyPodcastResults(flaws.KeyboardSerious.ContinueError(consts.KEY_BUFF_ERROR, err))
	}
	var readFiles, savedFiles int
	go GoDownloadError(ctxMedias, cancelMedias, errorStream, seriousStream)
	go GoStopKey(ctxMedias, cancelMedias, podcastData, keysEvents, keyStream)
	curStat := consts.CurStat{
		ReadFiles:   &readFiles,
		SavedFiles:  &savedFiles,
		MinDiskMbs:  progBounds.MinDisk,
		NetworkLoad: progBounds.LoadOption,
	}
	for i := 0; i < numWorkers; i++ {
		go GoDownloadAndSaveFiles(ctxMedias, curStat, mediaStream, doneStream, errorStream, httpMedia)
	}
	possibleFiles, varietyFiles, err := media.SaveDownloadedMedia(ctxMedias, podcastData, mediaStream, progBounds.LimitOption, httpMedia) // fix orders like above !!!
	close(mediaStream)
	disposeDownloaders(numWorkers, doneStream)
	err = firstErr(err, seriousStream)
	podcastResults := consts.PodcastResults{
		ReadFiles:     readFiles,
		SavedFiles:    savedFiles,
		PossibleFiles: possibleFiles,
		VarietyFiles:  varietyFiles,
		PodcastTime:   time.Since(startTime),
		Err:           nil,
	}
	resultsWithErr := dealWithErrors(err, ctxMedias.Err(), podcastData, podcastResults)
	return resultsWithErr
}

func dealWithErrors(err, ctxErr error, podcastData consts.PodcastData, podcastResults consts.PodcastResults) consts.PodcastResults {

	if ctxErr == context.Canceled {
		podcastResults.Err = flaws.SStop.StartError(podcastData.PodPath)
		return podcastResults
	}
	if err != nil {
		return misc.EmptyPodcastResults(err)
	}
	if ctxErr != nil {
		return misc.EmptyPodcastResults(ctxErr)
	}
	return podcastResults
}

func firstErr(err error, seriousStream <-chan error) error {
	if err != nil {
		return err
	}
	for i := 0; i < len(seriousStream); i++ {
		err := <-seriousStream
		if err != nil {
			return err
		}
	}
	return nil
}
