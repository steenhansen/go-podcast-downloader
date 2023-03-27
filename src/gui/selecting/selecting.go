package selecting

import (
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/gui/colors"
	"podcast-downloader/src/gui/state"
	"podcast-downloader/src/gui/values"

	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ChangeDebugMess(reduxState string) {
	if values.GUI_DEBUG {
		if values.An_exe_debug_error_message != "" {
			state.TheMediaWindow.DebugState.Text = values.An_exe_debug_error_message
		} else {
			state.TheMediaWindow.DebugState.Text = reduxState
		}
		state.TheMediaWindow.DebugState.Refresh()
	}
}

func debugState() *fyne.Container {
	debugLabel := widget.NewLabel("-current-redux-state-")
	horDebug := container.NewCenter(debugLabel)
	state.TheMediaWindow.DebugState = debugLabel
	if !values.GUI_DEBUG {
		horDebug.Hidden = true
	}
	return horDebug
}

func declareWhom() *fyne.Container {
	currentNameBtn := widget.NewButton("-default-is-replaced-on-choice-", func() {})
	currentNameBtn.Disable()
	horWhom := container.NewCenter(currentNameBtn)
	state.TheMediaWindow.WhomBox = horWhom
	state.TheMediaWindow.WhomButton = currentNameBtn
	return horWhom
}

func DirPodcasts(redrawWindow func(state.StateKind)) *fyne.Container {
	actionButtons := container.NewVBox()
	actionButtons.Add(debugState())

	actionButtons.Add(backButton(redrawWindow))
	actionButtons.Add(declareWhom())
	actionButtons.Add(titleOrFname(redrawWindow))
	actionButtons.Add(allNone(redrawWindow))
	actionButtons.Add(stopDownload(redrawWindow))

	thePodcasts := container.NewVBox()
	for _, podcastName := range state.TheMediaWindow.PodcastDirs {
		closedName := podcastName
		podcastButton := widget.NewButton(closedName, func() {
			redrawWindow(state.B_NOTHING_SELECTED)
			state.TheMediaWindow.CurPodcastDir = closedName
			mediaTitles, _ := state.GetRssFile(closedName)
			numberEpisodes := len(mediaTitles)
			state.TheMediaWindow.ChosenFnames = make([]bool, numberEpisodes)
			state.TheMediaWindow.ChosenTitles = make([]bool, numberEpisodes)
			state.SetRssType(closedName)
			redrawWindow(state.F_CHOOSEN_FILENAMES_NONE)
			downloadMess := `Downloading "` + closedName + `"`
			state.TheMediaWindow.WhomButton.SetText(downloadMess)
		})
		thePodcasts.Add(podcastButton)
	}
	state.TheMediaWindow.PodcastList = thePodcasts
	actionButtons.Add(thePodcasts)
	actionButtons.Add(startDownload(redrawWindow))
	leftScroll := container.NewVScroll(actionButtons)
	leftContainer := container.NewMax(leftScroll)
	return leftContainer
}

func backButton(redrawWindow func(state.StateKind)) *fyne.Container {
	var backIcon *fyne.StaticResource

	if values.UseDyanmicButtonIcons {
		backIcon = ButtonIcon("go-back.png")
	} else {
		backIcon = resourceGoBackPng
	}

	backButton := widget.NewButtonWithIcon("Back To Podcast List", backIcon, func() {
		redrawWindow(state.B_NOTHING_SELECTED)
	})
	horBack := container.NewHBox(backButton)
	state.TheMediaWindow.BackBox = horBack
	return horBack
}

func BB() *fyne.StaticResource {
	return resourceGoBackPng
}

// copy dos/tests_mocked_http/pressStop_m/pressStop_m_test.go
func stopDownload(redrawWindow func(state.StateKind)) *fyne.Container {
	var stopIcon *fyne.StaticResource
	if values.UseDyanmicButtonIcons {
		stopIcon = ButtonIcon("stop-downloading.png")
	} else {
		stopIcon = resourceStopDownloadingPng
	}

	stopDownloading := widget.NewButtonWithIcon("Stop", stopIcon, func() {
		state.TheMediaWindow.KeyStream <- "S" // sends stop signal to downloading terminal section
		state.TheMediaWindow.StopDownloadBox.Hidden = true
	})
	quitDownloading := container.NewHBox(stopDownloading)
	state.TheMediaWindow.StopDownloadBox = quitDownloading
	return quitDownloading
}

func allNone(redrawWindow func(state.StateKind)) *fyne.Container {
	var allIcon *fyne.StaticResource
	if values.UseDyanmicButtonIcons {
		allIcon = ButtonIcon("select-all.png")
	} else {
		allIcon = resourceSelectAllPng
	}

	selectAll := widget.NewButtonWithIcon("Select All Episodes", allIcon, func() {
		state.AllSelected(redrawWindow)
	})

	var noneIcon *fyne.StaticResource
	if values.UseDyanmicButtonIcons {
		noneIcon = ButtonIcon("select-none.png")
	} else {
		noneIcon = resourceSelectNonePng
	}

	selectNone := widget.NewButtonWithIcon("Select No Episodes", noneIcon, func() {
		state.NoneSelected(redrawWindow)
	})
	selectAllNone := container.NewHBox(selectAll, selectNone)
	state.TheMediaWindow.SelectAllOrNoneBox = selectAllNone
	return selectAllNone
}

func startDownload(redrawWindow func(state.StateKind)) *fyne.Container {
	var downloadResource *fyne.StaticResource
	if values.UseDyanmicButtonIcons {
		downloadResource = ButtonIcon("prog-icon.png")
	} else {
		downloadResource = resourceProgIconPng
	}
	downloadIcon := widget.NewIcon(downloadResource)
	indentedDownloadIcon := container.NewHBox(downloadIcon, downloadIcon)

	emptyDownloadBtn := widget.NewButton("", func() {
		redrawWindow(state.I_START_DOWNLOADING)
	})
	btnColor := canvas.NewRectangle(colors.TRANSLUCENT_BACKGROUND)
	btnText := canvas.NewText("Download Podcast Episodes", colors.BLACK_DOWNLOAD)
	btnText.Alignment = fyne.TextAlignCenter
	btnText.TextStyle = fyne.TextStyle{Bold: true}
	downloadBox := container.New(
		layout.NewMaxLayout(),
		emptyDownloadBtn,
		btnColor,
		btnText,
		indentedDownloadIcon,
	)
	state.TheMediaWindow.StartDownloadBox = downloadBox
	return downloadBox
}

func titleOrFname(redrawWindow func(state.StateKind)) *fyne.Container {
	fileTitle := state.FileOrTitle(redrawWindow)
	state.TheMediaWindow.TitleOrFnameBox = fileTitle
	return fileTitle
}

func RightTitles(redrawWindow func(state.StateKind)) *fyne.Container {
	mediaTitles, rssFiles := state.GetRssFile(state.TheMediaWindow.CurPodcastDir)
	numberEpisodes := len(mediaTitles)
	state.TheMediaWindow.EpisodeUrls = make([]string, numberEpisodes)
	copy(state.TheMediaWindow.EpisodeUrls, rssFiles)
	fileExt := media.FileExten(rssFiles[0])
	state.TheMediaWindow.PodcastFileExt = fileExt
	state.TheMediaWindow.EpisodeTitles = make([]string, numberEpisodes)

	state.TheMediaWindow.EpisodeErrors = make(map[string]string, numberEpisodes)

	copy(state.TheMediaWindow.EpisodeTitles, mediaTitles)
	var lengthChecks = func() int { return numberEpisodes }
	var createCheckbox = func() fyne.CanvasObject {
		newCheckbox := widget.NewCheck("building-check-box", nil)
		return container.NewPadded(newCheckbox)
	}
	fileDirectory := state.TheMediaWindow.ProgPath + "/" + state.TheMediaWindow.CurPodcastDir + "/"
	var updateCheckbox = func(checkId widget.ListItemID, item fyne.CanvasObject) {
		episodeTitle := state.TheMediaWindow.EpisodeTitles[checkId]
		safeTitle := media.TitleToName(episodeTitle)
		titleWithExt := safeTitle + "." + state.TheMediaWindow.PodcastFileExt
		nameOfFile := rss.NameOfFile(titleWithExt)
		filePath := fileDirectory + nameOfFile
		fixCheck := item.(*fyne.Container).Objects[0].(*widget.Check)
		fixCheck.Text = nameOfFile
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				fixCheck.Checked = state.TheMediaWindow.ChosenTitles[checkId]
				fixCheck.OnChanged = func(b bool) {
					state.TheMediaWindow.ChosenTitles[checkId] = b
					state.ManyOneNone(redrawWindow)
				}
			}
		} else {
			state.TheMediaWindow.ChosenTitles[checkId] = false
			fixCheck.Disable()
		}
		fixCheck.Refresh()
	}
	titleChecks := widget.NewList(lengthChecks, createCheckbox, updateCheckbox)
	scrollChecks := container.NewVScroll(titleChecks)
	rightContainer := container.NewMax(scrollChecks)
	return rightContainer
}

func RightFilenames(redrawWindow func(state.StateKind)) (*fyne.Container, bool) {
	mediaTitles, rssFiles := state.GetRssFile(state.TheMediaWindow.CurPodcastDir)
	numberEpisodes := len(mediaTitles)
	state.TheMediaWindow.EpisodeUrls = make([]string, numberEpisodes)
	copy(state.TheMediaWindow.EpisodeUrls, rssFiles)
	fileExt := media.FileExten(rssFiles[0])
	state.TheMediaWindow.PodcastFileExt = fileExt
	state.TheMediaWindow.EpisodeTitles = make([]string, numberEpisodes)
	state.TheMediaWindow.EpisodeErrors = make(map[string]string, numberEpisodes)
	copy(state.TheMediaWindow.EpisodeTitles, mediaTitles)
	var lengthChecks = func() int { return numberEpisodes }
	var createCheckbox = func() fyne.CanvasObject {
		newCheckbox := widget.NewCheck("building-check-box", nil)
		return container.NewPadded(newCheckbox)
	}
	fileDirectory := state.TheMediaWindow.ProgPath + "/" + state.TheMediaWindow.CurPodcastDir + "/"
	var updateCheckbox = func(checkId widget.ListItemID, item fyne.CanvasObject) {
		urlFilename := state.TheMediaWindow.EpisodeUrls[checkId]
		nameOfFile := rss.NameOfFile(urlFilename)
		filePath := fileDirectory + nameOfFile
		fixCheck := item.(*fyne.Container).Objects[0].(*widget.Check)
		fixCheck.Text = nameOfFile
		if _, err := os.Stat(filePath); err != nil {
			if os.IsNotExist(err) {
				fixCheck.Checked = state.TheMediaWindow.ChosenFnames[checkId]
				fixCheck.OnChanged = func(b bool) {
					state.TheMediaWindow.ChosenFnames[checkId] = b
					state.ManyOneNone(redrawWindow)
				}
			}
		} else {
			state.TheMediaWindow.ChosenFnames[checkId] = false
			fixCheck.Disable()
		}
		fixCheck.Refresh()
	}
	fnameChecks := widget.NewList(lengthChecks, createCheckbox, updateCheckbox)
	scrollChecks := container.NewVScroll(fnameChecks)
	rightContainer := container.NewMax(scrollChecks)
	// n.b. test to see if every filename is the same, like Timesuck's "default.mp3"
	allFnamesSame := false
	filenamesSet := make(map[string]string)
	for _, urlFilename := range state.TheMediaWindow.EpisodeUrls {
		shortName := rss.NameOfFile(urlFilename)
		filenamesSet[shortName] = shortName
	}
	if numberEpisodes > 1 && len(filenamesSet) == 1 {
		allFnamesSame = true
	}
	return rightContainer, allFnamesSame
}
