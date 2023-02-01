package processes

import (
	"context"
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

// simKeyStream <-chan string is so that a test can simulate stopping
func GoStopKey(ctx context.Context, cancel context.CancelFunc, dirName string, keysEvents <-chan keyboard.KeyEvent, simKeyStream <-chan string) {
	globals.Console.Note("Downloading '" + dirName + "' podcast, hit '" + consts.STOP_KEY_LOWER + "' to stop")
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
		case simKey := <-simKeyStream:
			globals.Console.Note("TESTING - downloading stopped by simulated key press of '" + simKey + "'")
			cancel()
			break keyboardCancel
		case <-ctx.Done():
			break keyboardCancel
		}
	}
}

func GoDownloadAndSaveFiles(ctx context.Context, curStat consts.CurStat, mediaStream <-chan consts.MediaEnclosure, doneStream chan<- bool, errorStream chan<- error) {
downloadCancel:
	for media := range mediaStream {
		select {
		case <-ctx.Done():
			break downloadCancel
		default:
			start := time.Now()
			globals.Console.Note(feed.ShowProgress(media, curStat.ReadFiles))
			err := rss.DownloadAndWriteFile(ctx, media.EnclosureUrl, media.EnclosurePath, curStat.MinDiskMbs)
			if err != nil {
				if ctx.Err() != context.Canceled {
					errorStream <- err
					globals.Faults.Note(media.EnclosureUrl, err)
					globals.Console.Note(feed.ShowError(media.EnclosureUrl))
				}
			} else {
				globals.Console.Note(feed.ShowSaved(curStat.SavedFiles, start, media.EnclosurePath))
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

func DownloadMedia(url string, podcastData consts.PodcastData, progBounds consts.ProgBounds, simKeyStream chan string) consts.PodcastResults {
	startTime := time.Now()
	numWorkers := misc.NumWorkers(progBounds.LoadOption)
	mediaStream := make(chan consts.MediaEnclosure)
	doneStream := make(chan bool, numWorkers)
	errorStream := make(chan error, numWorkers)
	seriousStream := make(chan error, numWorkers)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(consts.MAX_READ_FILE_TIME))
	defer cancel()
	keysEvents, err := keyboard.GetKeys(consts.KEY_BUFF_SIZE)
	if err != nil {
		return misc.EmptyPodcastResults(flaws.KeyboardSerious.ContinueError(consts.KEY_BUFF_ERROR, err))
	}
	var readFiles, savedFiles int
	go GoDownloadError(ctx, cancel, errorStream, seriousStream)
	go GoStopKey(ctx, cancel, podcastData.PodTitle, keysEvents, simKeyStream)
	curStat := consts.CurStat{
		ReadFiles:  &readFiles,
		SavedFiles: &savedFiles,
		MinDiskMbs: progBounds.MinDisk,
	}
	for i := 0; i < numWorkers; i++ {
		go GoDownloadAndSaveFiles(ctx, curStat, mediaStream, doneStream, errorStream)
	}
	possibleFiles, varietyFiles, err := media.SaveDownloadedMedia(ctx, podcastData, mediaStream, progBounds.LimitOption)
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
	resultsWithErr := dealWithErrors(err, ctx.Err(), podcastData, podcastResults)
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
