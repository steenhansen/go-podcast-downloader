package terminal

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/media"
	"github.com/steenhansen/go-podcast-downloader-console/src/podcasts"
	"github.com/steenhansen/go-podcast-downloader-console/src/processes"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func ShowNumberedChoices(progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	pdescs, thePodcasts, err := podcasts.AllPodcasts(progBounds.ProgPath)
	if err != nil {
		return "", err
	}
	if len(thePodcasts) == 0 {
		return "", flaws.NoPodcasts.StartError("add some podcasts feeds first")
	}
	choices, err := podcasts.PodChoices(progBounds.ProgPath, pdescs)
	if err != nil {
		return "", err
	}
	fmt.Print("\n", choices, " 'Q' or a number + enter: ")
	choice, err := podcasts.ChoosePod(progBounds.ProgPath, pdescs)
	if choice == 0 && err == nil {
		return "", nil // 'Q' entered to quit
	}
	if err != nil {
		return "", err
	}
	report, err := downloadAndReport(pdescs, thePodcasts, choice-1, progBounds, simKeyStream, mediaFix)
	if err != nil {
		return "", err
	}
	return report, nil
}

// go run ./ siberiantimes.com/ecology/rss/
func AddByUrl(url string, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	xml, files, lengths, err := podcasts.ReadUrl(url)
	if err != nil {
		return "", err
	}
	mediaTitle, err := rss.RssTitle(xml)
	if err != nil {
		return "", err
	}
	mediaPath, err := media.InitFolder(progBounds.ProgPath, mediaTitle, url)
	if err != nil {
		return "", err
	}
	podcastData := consts.PodcastData{
		MediaTitle: mediaTitle,
		MediaPath:  mediaPath,
		Medias:     files,
		Lengths:    lengths,
	}
	report, err := downloadReport(url, podcastData, progBounds, simKeyStream, mediaFix)
	return report, err
}

// go run ./ siberiantimes.com/ecology/rss/ Xecology
func AddByUrlAndName(url string, osArgs []string, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	_, files, lengths, err := podcasts.ReadUrl(url)
	if err != nil {
		return "", err
	}
	mediaTitle := feed.PodcastName(osArgs)

	mediaPath, err := media.InitFolder(progBounds.ProgPath, mediaTitle, url)
	if err != nil {
		return "", err
	}
	podcastData := consts.PodcastData{
		MediaTitle: mediaTitle,
		MediaPath:  mediaPath,
		Medias:     files,
		Lengths:    lengths,
	}
	report, err := downloadReport(url, podcastData, progBounds, simKeyStream, mediaFix)
	return report, err
}

func downloadAndReport(pdescs, feed []string, choice int, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	mediaTitle := pdescs[choice]
	url := feed[choice]
	podcastResults := podcasts.DownloadPodcast(mediaTitle, url, progBounds, simKeyStream, mediaFix)
	//fmt.Print("q*bert s is a non-stopping error podcastResults==", podcastResults)
	if podcastResults.Err != nil && errors.Is(podcastResults.Err, flaws.LowDisk) {

		return "", podcastResults.Err
	}
	report := doReport(podcastResults, string(url), mediaTitle)
	//fmt.Println("q*bert    report = ", report)
	return report, nil
}

func doReport(podcastResults consts.PodcastResults, url string, mediaTitle string) string {
	//readFiles := podcastResults.ReadFiles
	savedFiles := podcastResults.SavedFiles
	//possibleFiles := podcastResults.PossibleFiles               keep for tests
	varietyFiles := podcastResults.VarietyFiles
	podcastTime := podcastResults.PodcastTime
	rounded := podcastTime.Round(time.Second)
	//	report := fmt.Sprintf("\nGot %d new '%s' file(s) from %s into '%s' for a collection of %d out of a possible %d in %s",
	//savedFiles, varietyFiles, url, mediaTitle, readFiles, possibleFiles, rounded)
	report := fmt.Sprintf("\nAdded %d new '%s' file(s) from %s into '%s' in %s",
		savedFiles, varietyFiles, url, mediaTitle, rounded)
	//fmt.Println("jjjjjjjjjjjjjjjjjjjjj report = ", report)
	return report
}

func downloadReport(url string, podcastData consts.PodcastData, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	fmt.Print("\nADDING ", podcastData.MediaTitle, "\n\n")
	podcastResults := processes.DownloadMedia(url, podcastData, progBounds, simKeyStream, mediaFix)
	if podcastResults.Err != nil {
		return "", podcastResults.Err
	}
	report := doReport(podcastResults, url, podcastData.MediaTitle)
	return report, nil
}

func ReadByExistName(osArgs []string, progBounds consts.ProgBounds, simKeyStream chan string, mediaFix map[string]error) (string, error) {
	name := feed.PodcastName(osArgs)

	mediaPath, mediaTitle, err := podcasts.FindPodcastDirName(progBounds.ProgPath, name)
	if err != nil {
		return "", err
	}
	origin := mediaPath + "/" + consts.URL_OF_RSS
	url, err := os.ReadFile(origin)
	if err != nil {
		return "", err
	}
	strUrl := string(url)
	_, files, lengths, err := podcasts.ReadUrl(strUrl) // _ == unused xml
	if err != nil {
		return "", err
	}
	podcastData := consts.PodcastData{
		MediaTitle: mediaTitle,
		MediaPath:  mediaPath,
		Medias:     files,
		Lengths:    lengths,
	}
	report, err := downloadReport(string(url), podcastData, progBounds, simKeyStream, mediaFix)

	return report, err
}
