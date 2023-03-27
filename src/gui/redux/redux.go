package redux

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/initialize"
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/podcasts"
	"podcast-downloader/src/dos/terminal"
	"podcast-downloader/src/gui/backingup"
	"podcast-downloader/src/gui/selecting"
	"podcast-downloader/src/gui/state"
	"podcast-downloader/src/gui/values"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

func MenuUpdate(internetLoadChangable bool) {
	if state.TheMediaWindow.Internetload == consts.HIGH_LOAD {
		state.TheMediaWindow.Menu2High.Checked = true
		state.TheMediaWindow.Menu2Medium.Checked = false
		state.TheMediaWindow.Menu2Low.Checked = false
	} else if state.TheMediaWindow.Internetload == consts.MEDIUM_LOAD {
		state.TheMediaWindow.Menu2High.Checked = false
		state.TheMediaWindow.Menu2Medium.Checked = true
		state.TheMediaWindow.Menu2Low.Checked = false
	} else {
		state.TheMediaWindow.Menu2High.Checked = false
		state.TheMediaWindow.Menu2Medium.Checked = false
		state.TheMediaWindow.Menu2Low.Checked = true
	}
	state.TheMediaWindow.Menue1Add.Disabled = !internetLoadChangable
	state.TheMediaWindow.Menu2High.Disabled = !internetLoadChangable
	state.TheMediaWindow.Menu2Medium.Disabled = !internetLoadChangable
	state.TheMediaWindow.Menu2Low.Disabled = !internetLoadChangable
}

func WindowStart(fyneWindow fyne.Window) {
	state.TheMediaWindow.FyneWindow = fyneWindow
	initialize.AddNasa()
	state.TheMediaWindow.KeyStream = make(chan string)

	RedrawWindow(state.A_LOAD_DIRECTORIES) // to load directories
	RedrawWindow(state.B_NOTHING_SELECTED)
	MenuUpdate(true)
}

func aLoadDirectories() *fyne.Container {
	_, progBounds, _ := misc.InitProg()
	state.TheMediaWindow.ProgPath = progBounds.ProgPath
	podcastDirs, _, _, _ := podcasts.AllPodcasts(progBounds.ProgPath)
	state.TheMediaWindow.PodcastDirs = podcastDirs
	state.TheMediaWindow.LeftSide = selecting.DirPodcasts(RedrawWindow)
	rightContainer := container.NewVBox(widget.NewLabel("Adding ..."))
	selecting.ChangeDebugMess("A_LOAD_DIRECTORIES")
	state.TheMediaWindow.WhomBox.Hidden = true
	state.TheMediaWindow.TitleOrFnameBox.Hidden = true
	state.TheMediaWindow.BackBox.Hidden = true
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = true
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = false
	return rightContainer
}

func bNothingSelected() *fyne.Container {
	rightContainer := container.NewVBox(widget.NewLabel("Select a podcast on left side"))
	selecting.ChangeDebugMess("B_NOTHING_SELECTED")
	state.TheMediaWindow.WhomBox.Hidden = true
	state.TheMediaWindow.TitleOrFnameBox.Hidden = true
	state.TheMediaWindow.BackBox.Hidden = true
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = true
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = false
	return rightContainer
}

func dNoChosenTitles() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("C_CHOOSEN_TITLES_NONE")
	banner := "Chosen filenames are derived from episode titles"
	return noneChosenYet(banner)
}

func gNoFilenamesChosen() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("F_CHOOSEN_FILENAMES_NONE")
	banner := "Chosen filenames are actual files in RSS"
	return noneChosenYet(banner)
}

func noneChosenYet(banner string) (*fyne.Container, bool) {
	state.TheMediaWindow.DebugState.Refresh()

	state.TheMediaWindow.FyneWindow.SetTitle(banner)
	var rightContainer *fyne.Container
	var allFnamesSame = false
	if state.TheMediaWindow.ForceTitleOverFname {
		rightContainer = selecting.RightTitles(RedrawWindow)
	} else {
		rightContainer, allFnamesSame = selecting.RightFilenames(RedrawWindow)
	}
	state.TheMediaWindow.WhomBox.Hidden = false
	state.TheMediaWindow.TitleOrFnameBox.Hidden = false
	state.TheMediaWindow.BackBox.Hidden = false
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = false
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = true
	return rightContainer, allFnamesSame
}

func eOneTitleChosen() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("D_CHOOSEN_TITLES_ONE")
	banner := "The filename is derived from episode titles"
	return atLeastOne(banner)
}

func hOneFilenameChosen() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("G_CHOOSEN_FILENAMES_ONE")
	banner := "The filename is actual filename in RSS"
	return atLeastOne(banner)
}

func iManyFilenamesChosen() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("H_CHOOSEN_FILENAMES_MANY")
	banner := "The filenames are derived from episode titles"
	return atLeastOne(banner)
}

func fManyChosenTitles() (*fyne.Container, bool) {
	selecting.ChangeDebugMess("E_CHOOSEN_TITLES_MANY")
	banner := "The filenames are actual filename in RSS"
	return atLeastOne(banner)
}

func atLeastOne(banner string) (*fyne.Container, bool) {
	state.TheMediaWindow.DebugState.Refresh()

	state.TheMediaWindow.FyneWindow.SetTitle(banner)
	var rightContainer *fyne.Container
	var allFnamesSame = false
	if state.TheMediaWindow.ForceTitleOverFname {
		rightContainer = selecting.RightTitles(RedrawWindow)
	} else {
		rightContainer, allFnamesSame = selecting.RightFilenames(RedrawWindow)
	}
	state.TheMediaWindow.WhomBox.Hidden = false
	state.TheMediaWindow.TitleOrFnameBox.Hidden = false
	state.TheMediaWindow.BackBox.Hidden = false
	state.TheMediaWindow.StartDownloadBox.Hidden = false
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = false
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = true
	return rightContainer, allFnamesSame
}

func jStartDownloading() *fyne.Container {
	go backingup.CallTerminalDownloading(RedrawWindow, state.TheMediaWindow.KeyStream)
	progPath := state.TheMediaWindow.ProgPath
	podTitle := state.TheMediaWindow.CurPodcastDir
	rssUrl := state.TheMediaWindow.PodcastUrl
	forcingTitle := state.TheMediaWindow.ForceTitleOverFname
	media.ReSaveFolder(forcingTitle, progPath, podTitle, rssUrl)
	selecting.ChangeDebugMess("I_START_DOWNLOADING")
	rightContainer := container.NewVBox(widget.NewLabel("Downloading ..."))
	return rightContainer
}

func kShowDownloading(redrawWindow func(state.StateKind)) *fyne.Container {
	var rightContainer *fyne.Container
	var remainingCount, totalCount int
	if state.TheMediaWindow.ForceTitleOverFname {
		rightContainer, remainingCount, totalCount = backingup.LabelTitles()
	} else {
		rightContainer, remainingCount, totalCount = backingup.LabelFilenames()
	}
	remCount := strconv.Itoa(remainingCount)
	totCount := strconv.Itoa(totalCount)
	downloadingTitle := " " + remCount + "/" + totCount + "   " + state.TheMediaWindow.SpinChar
	state.TheMediaWindow.FyneWindow.SetTitle(downloadingTitle)
	selecting.ChangeDebugMess("J_ARE_DOWNLOADING")
	state.TheMediaWindow.WhomBox.Hidden = false
	state.TheMediaWindow.TitleOrFnameBox.Hidden = true
	state.TheMediaWindow.BackBox.Hidden = true
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = true
	state.TheMediaWindow.StopDownloadBox.Hidden = false
	state.TheMediaWindow.PodcastList.Hidden = true
	return rightContainer
}

func lStopping() *fyne.Container {
	selecting.ChangeDebugMess("K_STOPPING")
	state.TheMediaWindow.WhomBox.Hidden = false
	state.TheMediaWindow.TitleOrFnameBox.Hidden = true
	state.TheMediaWindow.BackBox.Hidden = true
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = true
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = true
	rightContainer := container.NewVBox(widget.NewLabel("Reporting ..."))
	return rightContainer
}

func mReporting() *fyne.Container {
	selecting.ChangeDebugMess("L_REPORTING")
	state.TheMediaWindow.FyneWindow.SetTitle(values.WINDOW_DEFAULT_TITLE)
	podcastResults := state.TheMediaWindow.DownloadResults
	urlStr := state.TheMediaWindow.PodcastUrl
	podTitle := state.TheMediaWindow.CurPodcastDir
	rightContainer := container.NewVBox()
	podcastReport := terminal.DoReport(podcastResults, urlStr, podTitle)
	resultsLabel := widget.NewLabel(podcastReport)
	rightContainer.Add(resultsLabel)
	if len(state.TheMediaWindow.EpisodeErrors) > 0 {
		rightContainer.Add(widget.NewLabel("ERRORS"))
		for _, episodeFname := range state.TheMediaWindow.EpisodeErrors {
			errorLabel := widget.NewLabel(episodeFname)
			rightContainer.Add(errorLabel)
		}
	}
	state.TheMediaWindow.WhomBox.Hidden = false
	state.TheMediaWindow.TitleOrFnameBox.Hidden = true
	state.TheMediaWindow.BackBox.Hidden = false
	state.TheMediaWindow.StartDownloadBox.Hidden = true
	state.TheMediaWindow.SelectAllOrNoneBox.Hidden = true
	state.TheMediaWindow.StopDownloadBox.Hidden = true
	state.TheMediaWindow.PodcastList.Hidden = true
	return rightContainer
}

func RedrawWindow(newState state.StateKind) {
	var rightContainer *fyne.Container
	var allFnamesSame bool
	switch newState {
	case state.A_LOAD_DIRECTORIES:
		MenuUpdate(true)
		rightContainer = aLoadDirectories()
	case state.B_NOTHING_SELECTED:
		MenuUpdate(true)
		rightContainer = bNothingSelected()
	case state.C_CHOOSEN_TITLES_NONE:
		MenuUpdate(true)
		rightContainer, _ = dNoChosenTitles()
	case state.D_CHOOSEN_TITLES_ONE:
		MenuUpdate(true)
		rightContainer, _ = eOneTitleChosen()
	case state.E_CHOOSEN_TITLES_MANY:
		MenuUpdate(true)
		rightContainer, _ = fManyChosenTitles()
	case state.F_CHOOSEN_FILENAMES_NONE:
		MenuUpdate(true)
		rightContainer, allFnamesSame = gNoFilenamesChosen()
		if allFnamesSame {
			state.TheMediaWindow.TitleOrFnameRd.SetSelected(state.TITLE_TYPE)
			return
		}
	case state.G_CHOOSEN_FILENAMES_ONE:
		MenuUpdate(true)
		rightContainer, allFnamesSame = hOneFilenameChosen()
		if allFnamesSame {
			state.TheMediaWindow.TitleOrFnameRd.SetSelected(state.TITLE_TYPE)
			return
		}
	case state.H_CHOOSEN_FILENAMES_MANY:
		MenuUpdate(true)
		rightContainer, allFnamesSame = iManyFilenamesChosen()
		if allFnamesSame {
			state.TheMediaWindow.TitleOrFnameRd.SetSelected(state.TITLE_TYPE)
			return
		}
	case state.I_START_DOWNLOADING:
		MenuUpdate(false)
		rightContainer = jStartDownloading()
	case state.J_ARE_DOWNLOADING:
		MenuUpdate(false)
		rightContainer = kShowDownloading(RedrawWindow)
	case state.K_STOPPING:
		MenuUpdate(true)
		rightContainer = lStopping()
	case state.L_REPORTING:
		MenuUpdate(true)
		rightContainer = mReporting()
	default:
		fmt.Println("ERROR Unknown redux = ", newState)
	}
	leftContainer := state.TheMediaWindow.LeftSide
	leftSplitRight2 := NewPointerHSplit(leftContainer, rightContainer)
	state.TheMediaWindow.FyneWindow.SetContent(leftSplitRight2)

	if newState == state.I_START_DOWNLOADING {
		RedrawWindow(state.J_ARE_DOWNLOADING)
	} else if newState == state.K_STOPPING {
		RedrawWindow(state.L_REPORTING)
	}
}

type PointerHSplit struct {
	container.Split
}

func NewPointerHSplit(leading, trailing fyne.CanvasObject) *PointerHSplit {
	hsplit := &PointerHSplit{}
	hsplit.ExtendBaseWidget(hsplit) // n.b. no other way of doing busy mouse pointer
	hsplit.Leading = leading
	hsplit.Trailing = trailing
	hsplit.Horizontal = true
	return hsplit
}

type BusyCursor int

func (d BusyCursor) Image() (image.Image, int, int) {
	var busyCursor image.Image
	if values.UseDyanmicButtonIcons {
		curDir := misc.CurDir()
		filePath := curDir + "/src/gui/images/busy-cursor.png"
		iconFile, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		r := bufio.NewReader(iconFile)
		busyCursor, _ = png.Decode(r)
	} else {
		busyCursor, _, _ = image.Decode(bytes.NewReader(resourceBusyCursorPng.StaticContent))
	}
	return busyCursor, 0, 0
}

func (hs *PointerHSplit) Cursor() desktop.Cursor {
	if state.TheMediaWindow.BusyCursor {
		var busyCursor BusyCursor
		return busyCursor
	}
	return desktop.DefaultCursor
}
