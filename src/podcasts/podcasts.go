package podcasts

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/feed"
	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/processes"
	"github.com/steenhansen/go-podcast-downloader-console/src/rss"
)

func FindPodcastDirName(path, name string) (string, string, error) {
	directory, err := os.Open(path)
	if err != nil {
		return "", "", err
	}
	defer directory.Close()

	dirs, err := directory.Readdir(0)
	if err != nil {
		return "", "", err
	}
	lowername := strings.ToLower(name)
	for _, dir := range dirs {
		if !dir.Mode().IsRegular() {
			MediaTitle := dir.Name()
			lowerdir := strings.ToLower(MediaTitle)
			if lowerdir == lowername {
				mediaPath := path + "/" + dir.Name()
				return mediaPath, MediaTitle, nil
			}
		}
	}
	return "", "", flaws.BadChoice.StartError(name)
}

func ReadUrl(url string) ([]byte, []string, []int, error) {
	xml, err := feed.ReadRss(url)
	if err != nil {
		return nil, nil, nil, err
	}
	title, err := rss.RssTitle(xml)
	if title == "" {
		xmlStr := string(xml[0:100])
		return nil, nil, nil, flaws.InvalidXML.StartError(xmlStr)
	} else if err != nil {
		return nil, nil, nil, err
	}
	files, lengths, err := rss.RssItems(xml)
	if err != nil {
		return nil, nil, nil, err
	}
	return xml, files, lengths, nil
}

func DownloadPodcast(mediaTitle, url string, progBounds consts.ProgBounds, simKeyStream chan string) consts.PodcastResults {
	if feed.IsUrl(url) {
		_, files, lengths, err := ReadUrl(url)
		if err != nil {
			return misc.EmptyPodcastResults(err)
		}
		mediaPath := progBounds.ProgPath + "/" + mediaTitle
		podcastData := consts.PodcastData{
			MediaTitle: mediaTitle,
			MediaPath:  mediaPath,
			Medias:     files,
			Lengths:    lengths,
		}
		podcastResults := processes.DownloadMedia(url, podcastData, progBounds, simKeyStream)
		return podcastResults
	}
	return misc.EmptyPodcastResults(flaws.InvalidRssURL.StartError(url))
}

// path =>dirname
func PodChoices(path string, pdescs []string) (choices string, err error) {
	var sizedStr string
	for i, dirdesc := range pdescs {
		count, size, ftypes, err := countFiles(path, dirdesc)
		if ftypes == "" && strings.Contains(dirdesc, "[") {
			fparts := strings.Split(dirdesc, "[")
			extension := fparts[1]
			ftypes = extension[0:3]
		}
		if err != nil {
			return "", err
		}
		if size < consts.GB_BYTES {
			mbs := size / consts.MB_BYTES
			sizedStr = fmt.Sprintf("%dMB", mbs)
		} else {
			gbs := float64(size) / float64(consts.GB_BYTES)
			sizedStr = fmt.Sprintf("%.2fGB", gbs)
		}
		choices += fmt.Sprintf("%2d | %12s |%4d files |%7s | %s\n", i+1, ftypes, count-1, sizedStr, dirdesc)
	}
	return choices, nil
}

// path =>dirname
func ChoosePod(path string, pdescs []string, getMenuChoice consts.ReadLineFunc) (choice int, err error) {
	input := getMenuChoice()
	text := strings.Trim(input, "\r\n")
	if text == "q" || text == "Q" {
		return 0, nil
	}
	choice, _ = strconv.Atoi(text)
	if choice < 1 || choice > len(pdescs) {
		return 0, flaws.BadChoice.StartError(text)
	}
	return choice, nil
}

// path =>dirname
func countFiles(path, dirName string) (count int, size int64, fTypes string, err error) {
	varieties := misc.VarietiesSet{}
	dirPath := path + "/" + dirName
	directory, err := os.Open(dirPath)
	if err != nil {
		return 0, 0, "", err
	}
	defer directory.Close()

	dirs, err := directory.Readdir(0)
	if err != nil {
		return 0, 0, "", err
	}

	for _, afile := range dirs {
		if afile.Mode().IsRegular() {
			varieties.AddVariety(afile.Name())
			size = size + afile.Size()
			count++
		}
	}
	fTypes = strings.TrimSpace(fTypes)
	fTypes = varieties.VarietiesString(" ")
	return count, size, fTypes, nil
}

// path =>dirname
func AllPodcasts(path string) ([]string, []string, error) {
	directory, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer directory.Close()

	dirs, err := directory.Readdir(0)
	if err != nil {
		return nil, nil, err
	}
	pdescs := make([]string, 0)
	feed := make([]string, 0)
	for _, dir := range dirs {
		if !dir.Mode().IsRegular() {
			dirname := dir.Name()
			if dirname != consts.SOURCE_FOLDER {
				origin := path + "/" + dir.Name() + "/" + consts.URL_OF_RSS
				rss, err := os.ReadFile(origin)
				if err == nil {
					pdescs = append(pdescs, dirname)
					feed = append(feed, string(rss))
				}

			}
		}
	}
	return pdescs, feed, nil
}
