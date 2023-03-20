package processes

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/feed"
	"podcast-downloader/src/dos/flaws"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/dos/stop"

	"github.com/eiannone/keyboard"
)

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

func Go_downloadMedia(ctx context.Context, curStat models.CurStat, mediaStream <-chan models.MediaEnclosure,
	errorStream chan<- models.MediaError, httpMedia models.HttpFn, waitGroup *sync.WaitGroup, afterDownloadEpisodeEvent func(string)) {
	globals.WaitCount.Adding()
	misc.ChannelLog("\t\t\t\t Go_downloadMedia START " + strconv.Itoa(globals.WaitCount.Current()))
	waitGroup.Add(1)
	for newMedia := range mediaStream {
		misc.SleepTime(curStat.NetworkLoad)
		start := time.Now()
		globals.Console.Note(feed.ShowProgress(newMedia, &readFiles))
		writtenBytes, err := rss.DownloadAndWriteFile(ctx, newMedia.EnclosureUrl, newMedia.EnclosurePath, curStat.MinDiskMbs, httpMedia)
		if ctx.Err() != context.Canceled && err != nil {
			mediaError := models.MediaError{
				EnclosureUrl:  newMedia.EnclosureUrl,
				EnclosurePath: newMedia.EnclosurePath,
				OrgErr:        err,
			}
			errorStream <- mediaError
		} else if ctx.Err() == nil && writtenBytes > 0 {
			afterDownloadEpisodeEvent(newMedia.EnclosureUrl)
			misc.ChannelLog("\t\t\t\t Go_downloadMedia SAVED " + newMedia.EnclosurePath)
			globals.Console.Note(feed.ShowSaved(&savedFiles, start, newMedia.EnclosurePath))
			globals.Console.Note(feed.ShowSizeError(newMedia.EnclosureSize, writtenBytes))
		}

	}
	waitGroup.Done()
	globals.WaitCount.Subtracting()
	misc.ChannelLog("\t\t\t\t Go_downloadMedia END " + strconv.Itoa(globals.WaitCount.Current()))
}

func createChannels(podcastData models.PodcastData, progBounds models.ProgBounds, keyStreamTest chan string,
	httpMedia models.HttpFn, ctxMedias context.Context, cancelMedias context.CancelFunc, afterDownloadEpisodeEvent func(string), downloadEpisodeErrorEvent func(string)) {
	numWorkers := misc.NumWorkers(progBounds.LoadOption)
	mediaStream = make(chan models.MediaEnclosure)
	errorStream = make(chan models.MediaError, numWorkers)
	seriousStream = make(chan error, numWorkers) // if run out of disk space, this stream could get 15 serious errors
	signalEndSerious = make(chan bool)
	signalEndStop = make(chan bool)
	KeyEventsReal, err := keyboard.GetKeys(consts.KEY_BUFF_SIZE)
	if err != nil {
		mediaError := models.MediaError{
			EnclosureUrl:  "",
			EnclosurePath: "",
			OrgErr:        err,
		}
		errorStream <- mediaError
	}
	go stop.Go_ctxDone(ctxMedias)
	go stop.Go_seriousError(ctxMedias, cancelMedias, errorStream, seriousStream, signalEndSerious, downloadEpisodeErrorEvent)
	go stop.Go_stopKey(cancelMedias, KeyEventsReal, keyStreamTest, signalEndStop, afterDownloadEpisodeEvent)

	startDownloading(cancelMedias, podcastData)
	curStat := models.CurStat{
		MinDiskMbs:  progBounds.MinDisk,
		NetworkLoad: progBounds.LoadOption,
	}
	for i := 0; i < numWorkers; i++ {
		go Go_downloadMedia(ctxMedias, curStat, mediaStream, errorStream, httpMedia, &waitGroup, afterDownloadEpisodeEvent)
	}
}

var readFiles, savedFiles int // Count number of files delt with

var waitGroup sync.WaitGroup // Controls Go_downloadMedia()

var mediaStream chan models.MediaEnclosure // Queue of media files to be downloaded
var errorStream chan models.MediaError     // All errors

var seriousStream chan error   // Out of Disk Space error
var signalEndSerious chan bool // Leave Go_seriousError()
var signalEndStop chan bool    // Leave Go_stopKey()

func BackupPodcast(url string, podcastData models.PodcastData, progBounds models.ProgBounds, keyStreamTest chan string,
	httpMedia models.HttpFn, afterDownloadEpisodeEvent func(string), downloadEpisodeErrorEvent func(string)) models.PodcastResults {
	misc.ChannelLog("BackupPodcast START")
	savedFiles = 0
	startTime := time.Now()
	timeOut := misc.FileTimeout(globals.MediaMaxReadFileTime)
	ctxMedias, cancelMedias := context.WithTimeout(context.Background(), timeOut)
	defer cancelMedias()
	createChannels(podcastData, progBounds, keyStreamTest, httpMedia, ctxMedias, cancelMedias, afterDownloadEpisodeEvent, downloadEpisodeErrorEvent)
	possibleFiles, osFileErr := media.Go_deriveFilenames(ctxMedias, podcastData, mediaStream, progBounds.LimitOption, httpMedia)
	misc.ChannelLog("BackupPodcast Go_deriveFilenames done")

	close(mediaStream)
	misc.ChannelLog(" Go_downloadMedia close(mediaStream)")
	waitGroup.Wait()

	misc.ChannelLog("BackupPodcast Go_downloadMedia done")
	signalEndSerious <- false
	signalEndStop <- true
	misc.ChannelLog("BackupPodcast all channels done")
	wasCanceled := false
	seriousError := firstErr(osFileErr, seriousStream)
	misc.ChannelLog("BackupPodcast serious error")

	if ctxMedias.Err() == context.Canceled {
		wasCanceled = true
	} else if ctxMedias.Err() != nil {
		seriousError = ctxMedias.Err()
		if ctxMedias.Err().Error() == consts.CONTEXT_DEAD_EXCEEDED {
			exceedTimeout := fmt.Sprintf("duration: %s", globals.MediaMaxReadFileTime)
			seriousError = flaws.TimeoutStop.MakeFlaw(exceedTimeout)
		}
	}
	if osFileErr != nil && osFileErr != context.Canceled {
		seriousError = osFileErr
	}
	misc.ChannelLog("BackupPodcast models.PodcastResults")
	podcastResults := models.PodcastResults{
		ReadFiles:     readFiles,
		SavedFiles:    savedFiles,
		PossibleFiles: possibleFiles,
		PodcastTime:   time.Since(startTime),
		WasCanceled:   wasCanceled,
		SeriousError:  seriousError,
	}
	misc.ChannelLog("BackupPodcast END")
	return podcastResults
}
