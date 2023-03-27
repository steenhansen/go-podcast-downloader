package backingup

import (
	"os"
	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/podcasts"
	"podcast-downloader/src/dos/processes"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/gui/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func LabelTitles() (*fyne.Container, int, int) {
	notExistingChosenTitles := make([]string, 0)
	for episodeIndex, wasChosen := range state.TheMediaWindow.ChosenTitles {
		if wasChosen {
			episodeTitle := state.TheMediaWindow.EpisodeTitles[episodeIndex]
			safeTitle := media.TitleToName(episodeTitle)
			titleWithExt := safeTitle + "." + state.TheMediaWindow.PodcastFileExt
			nameOfFile := rss.NameOfFile(titleWithExt)
			titleFname := state.TheMediaWindow.ProgPath + "/" + state.TheMediaWindow.CurPodcastDir + "/" + nameOfFile
			if _, err := os.Stat(titleFname); err != nil {
				if os.IsNotExist(err) {
					notExistingChosenTitles = append(notExistingChosenTitles, nameOfFile)
				}
			}
		}
	}
	titlesLeftDownload := len(notExistingChosenTitles)
	var lengthChecks = func() int { return titlesLeftDownload }
	var createLabel = func() fyne.CanvasObject {
		newLabel := widget.NewLabel("building-label-title")
		return container.NewPadded(newLabel)
	}
	var updateLabel = func(labelId widget.ListItemID, item fyne.CanvasObject) {
		nameOfFile := notExistingChosenTitles[labelId]
		fixLabel := item.(*fyne.Container).Objects[0].(*widget.Label)
		fixLabel.Text = nameOfFile
		_, containsKey := state.TheMediaWindow.EpisodeErrors[nameOfFile]
		if containsKey {
			fixLabel.Alignment = fyne.TextAlignTrailing
			fixLabel.TextStyle = fyne.TextStyle{Italic: true}
			fixLabel.Text = "ERROR - " + nameOfFile
		} else {
			fixLabel.Text = nameOfFile
		}
		fixLabel.Refresh()
	}
	titleLabels := widget.NewList(lengthChecks, createLabel, updateLabel)
	rightContainer := container.NewMax(titleLabels)
	numberChosen := len(state.TheMediaWindow.ChosenTitles)
	return rightContainer, titlesLeftDownload, numberChosen
}

func LabelFilenames() (*fyne.Container, int, int) {
	notExistingChosenFnames := make([]string, 0)
	for episodeIndex, wasChosen := range state.TheMediaWindow.ChosenFnames {
		if wasChosen {
			urlFilename := state.TheMediaWindow.EpisodeUrls[episodeIndex]
			shortName := rss.NameOfFile(urlFilename)
			titleFname := state.TheMediaWindow.ProgPath + "/" + state.TheMediaWindow.CurPodcastDir + "/" + shortName
			if _, err := os.Stat(titleFname); err != nil {
				if os.IsNotExist(err) {
					notExistingChosenFnames = append(notExistingChosenFnames, shortName)
				}
			}
		}
	}
	fnamesLeftDownload := len(notExistingChosenFnames)
	var lengthChecks = func() int { return fnamesLeftDownload }
	var createLabel = func() fyne.CanvasObject {
		newLabel := widget.NewLabel("building-label-fname")
		return container.NewPadded(newLabel)
	}
	var updateLabel = func(labelId widget.ListItemID, item fyne.CanvasObject) {
		nameOfFile := notExistingChosenFnames[labelId]
		fixLabel := item.(*fyne.Container).Objects[0].(*widget.Label)
		fixLabel.Text = nameOfFile
		_, containsKey := state.TheMediaWindow.EpisodeErrors[nameOfFile]
		if containsKey {
			fixLabel.Alignment = fyne.TextAlignTrailing
			fixLabel.TextStyle = fyne.TextStyle{Italic: true}
			fixLabel.Text = "ERROR - " + nameOfFile
		} else {
			fixLabel.Text = nameOfFile
		}
		fixLabel.Refresh()
	}
	titleLabels := widget.NewList(lengthChecks, createLabel, updateLabel)
	rightContainer := container.NewMax(titleLabels)
	numberChosen := len(state.TheMediaWindow.ChosenFnames)
	return rightContainer, fnamesLeftDownload, numberChosen
}

func CallTerminalDownloading(redrawWindow func(state.StateKind), keyStreamTest chan string) {
	CurPodcastDir := state.TheMediaWindow.CurPodcastDir
	mediaPath := state.TheMediaWindow.ProgPath + "/" + CurPodcastDir + "/"
	rssFilePath := state.TheMediaWindow.ProgPath + "/" + CurPodcastDir + "/" + consts.URL_OF_RSS_FN
	_, urlStr, _ := podcasts.IsForceTitle(rssFilePath)

	nMediaTitle := make([]string, 0)
	nRssSizes := make([]int, 0)
	nRssFiles := make([]string, 0)

	if state.TheMediaWindow.ForceTitleOverFname {
		for ind, isSel := range state.TheMediaWindow.ChosenTitles {
			if isSel {
				curUrl := state.TheMediaWindow.EpisodeUrls[ind]
				curTitle := state.TheMediaWindow.EpisodeTitles[ind]
				nMediaTitle = append(nMediaTitle, curTitle)
				nRssSizes = append(nRssSizes, 0)
				nRssFiles = append(nRssFiles, curUrl)
			}
		}
	} else {
		for ind, isSel := range state.TheMediaWindow.ChosenFnames {
			if isSel {
				curUrl := state.TheMediaWindow.EpisodeUrls[ind]
				curTitle := state.TheMediaWindow.EpisodeTitles[ind]
				nMediaTitle = append(nMediaTitle, curTitle)
				nRssSizes = append(nRssSizes, 0)
				nRssFiles = append(nRssFiles, curUrl)
			}
		}
	}

	_, progBounds, _ := misc.InitProg()
	podcastData := models.PodcastData{
		PodTitle:  CurPodcastDir,
		PodPath:   mediaPath,
		PodUrls:   nRssFiles,
		PodSizes:  nRssSizes,
		PodTitles: nMediaTitle,
	}

	progBounds.LoadOption = state.TheMediaWindow.Internetload
	afterDownloadEpisodeEvent := func(spinningSlashes string) {
		if len(spinningSlashes) == 0 || spinningSlashes[0] == consts.ASCII_CARRIAGE_RETURN {
			state.TheMediaWindow.FyneWindow.SetTitle(spinningSlashes)
		} else {
			state.TheMediaWindow.SpinChar = spinningSlashes
			redrawWindow(state.J_ARE_DOWNLOADING) // redraw window since files have changed
		}
	}

	downloadEpisodeErrorEvent := func(episodeFname string) {
		state.TheMediaWindow.EpisodeErrors[episodeFname] = episodeFname
	}

	globals.ForceTitle = state.TheMediaWindow.ForceTitleOverFname
	podcastResults := processes.BackupPodcast(urlStr, podcastData, progBounds, keyStreamTest, rss.HttpReal, afterDownloadEpisodeEvent, downloadEpisodeErrorEvent)
	state.TheMediaWindow.DownloadResults = podcastResults
	state.TheMediaWindow.PodcastUrl = urlStr
	redrawWindow(state.K_STOPPING)
}
