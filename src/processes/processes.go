package processes

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/steenhansen/go-podcast-downloader/src/consts"
	"github.com/steenhansen/go-podcast-downloader/src/feed"
	"github.com/steenhansen/go-podcast-downloader/src/globals"
	"github.com/steenhansen/go-podcast-downloader/src/media"
	"github.com/steenhansen/go-podcast-downloader/src/misc"
	"github.com/steenhansen/go-podcast-downloader/src/models"
	"github.com/steenhansen/go-podcast-downloader/src/rss"
	"github.com/steenhansen/go-podcast-downloader/src/stop"
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
	errorStream chan<- models.MediaError, httpMedia models.HttpFn, waitGroup *sync.WaitGroup) {
	globals.WaitCountDebug++
	misc.ChannelLog("\t\t\t\t Go_downloadMedia START " + strconv.Itoa(globals.WaitCountDebug))
	waitGroup.Add(1)
	for newMedia := range mediaStream {
		misc.SleepTime(curStat.NetworkLoad)
		start := time.Now()
		globals.Console.Note(feed.ShowProgress(newMedia, &readFiles))
		writtenBytes, err := rss.DownloadAndWriteFile(ctx, newMedia.EnclosureUrl, newMedia.EnclosurePath, curStat.MinDiskMbs, httpMedia)
		misc.ChannelLog("\t\t\t\t Go_downloadMedia SAVED " + newMedia.EnclosurePath)
		if ctx.Err() != context.Canceled && err != nil {
			mediaError := models.MediaError{
				EnclosureUrl:  newMedia.EnclosureUrl,
				EnclosurePath: newMedia.EnclosurePath,
				OrgErr:        err,
			}
			errorStream <- mediaError
		} else {
			if ctx.Err() == nil && writtenBytes > 0 && newMedia.EnclosureSize != writtenBytes {
				globals.Console.Note(feed.ShowSaved(&savedFiles, start, newMedia.EnclosurePath))
				globals.Console.Note(feed.ShowSizeError(newMedia.EnclosureSize, writtenBytes))
			}
		}
	}
	waitGroup.Done()
	misc.ChannelLog("\t\t\t\t Go_downloadMedia END " + strconv.Itoa(globals.WaitCountDebug))
	globals.WaitCountDebug--
}

func createChannels(podcastData models.PodcastData, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn, ctxMedias context.Context, cancelMedias context.CancelFunc) {
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
	go stop.Go_seriousError(ctxMedias, cancelMedias, errorStream, seriousStream, signalEndSerious)
	go stop.Go_stopKey(cancelMedias, KeyEventsReal, keyStream, signalEndStop)
	startDownloading(cancelMedias, podcastData)
	curStat := models.CurStat{
		MinDiskMbs:  progBounds.MinDisk,
		NetworkLoad: progBounds.LoadOption,
	}
	for i := 0; i < numWorkers; i++ {
		go Go_downloadMedia(ctxMedias, curStat, mediaStream, errorStream, httpMedia, &waitGroup)
	}
}

var readFiles, savedFiles int // Count number of files delt with

var waitGroup sync.WaitGroup // Controls Go_downloadMedia()

var mediaStream chan models.MediaEnclosure // Queue of media files to be downloaded
var errorStream chan models.MediaError     // All errors

var seriousStream chan error   // Out of Disk Space error
var signalEndSerious chan bool // Leave Go_seriousError()
var signalEndStop chan bool    // Leave Go_stopKey()

func BackupPodcast(url string, podcastData models.PodcastData, progBounds models.ProgBounds, keyStream chan string, httpMedia models.HttpFn) models.PodcastResults {
	misc.ChannelLog("BackupPodcast START")
	startTime := time.Now()
	timeOut := misc.FileTimeout(consts.MEDIA_MAX_READ_FILE_TIME)
	ctxMedias, cancelMedias := context.WithTimeout(context.Background(), timeOut)
	defer cancelMedias()
	createChannels(podcastData, progBounds, keyStream, httpMedia, ctxMedias, cancelMedias)
	possibleFiles, osFileErr := media.Go_deriveFilenames(ctxMedias, podcastData, mediaStream, progBounds.LimitOption, httpMedia)
	misc.ChannelLog("BackupPodcast Go_deriveFilenames done")

	close(mediaStream)
	misc.ChannelLog("BackupPodcast Go_downloadMedia close(mediaStream)")
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
