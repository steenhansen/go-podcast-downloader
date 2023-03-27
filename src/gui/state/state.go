package state

import (
	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/models"
	"podcast-downloader/src/dos/podcasts"
	"podcast-downloader/src/dos/rss"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func GetRssFile(dirName string) ([]string, []string) {
	TheMediaWindow.BusyCursor = true
	rssFilePath := TheMediaWindow.ProgPath + "/" + dirName + "/" + consts.URL_OF_RSS_FN
	_, urlStr, _ := podcasts.IsForceTitle(rssFilePath)
	TheMediaWindow.PodcastUrl = urlStr
	_, mediaTitles, rssFiles, _, _ := podcasts.ReadRssUrl(urlStr, rss.HttpReal)
	TheMediaWindow.BusyCursor = false
	return mediaTitles, rssFiles
}

type StateKind int

const (
	A_LOAD_DIRECTORIES StateKind = iota + 1
	B_NOTHING_SELECTED

	C_CHOOSEN_TITLES_NONE
	D_CHOOSEN_TITLES_ONE
	E_CHOOSEN_TITLES_MANY

	F_CHOOSEN_FILENAMES_NONE
	G_CHOOSEN_FILENAMES_ONE
	H_CHOOSEN_FILENAMES_MANY

	I_START_DOWNLOADING
	J_ARE_DOWNLOADING
	K_STOPPING
	L_REPORTING
)

type MediaWindow struct {
	PodcastDirs   []string
	CurPodcastDir string
	FyneWindow    fyne.Window
	ProgPath      string

	EpisodeUrls   []string
	EpisodeTitles []string

	EpisodeErrors map[string]string

	ForceTitleOverFname bool
	PodcastFileExt      string

	DownloadResults models.PodcastResults
	PodcastUrl      string

	LeftSide    *fyne.Container
	PodcastList *fyne.Container

	DebugState *widget.Label
	WhomBox    *fyne.Container
	WhomButton *widget.Button

	TitleOrFnameBox *fyne.Container
	TitleOrFnameRd  *widget.RadioGroup

	ChosenTitles []bool
	ChosenFnames []bool

	BackBox            *fyne.Container
	StartDownloadBox   *fyne.Container
	SelectAllOrNoneBox *fyne.Container
	StopDownloadBox    *fyne.Container

	KeyStream    chan string
	SpinChar     string
	Internetload string

	Menue1Add   *fyne.MenuItem
	Menu2High   *fyne.MenuItem
	Menu2Medium *fyne.MenuItem
	Menu2Low    *fyne.MenuItem

	BusyCursor bool
}

var TheMediaWindow = MediaWindow{
	PodcastDirs:   nil,
	CurPodcastDir: "",

	FyneWindow: nil,

	ProgPath:            "",
	EpisodeUrls:         nil,
	EpisodeTitles:       nil,
	EpisodeErrors:       nil,
	ForceTitleOverFname: false,

	PodcastFileExt:  "mp3",
	DownloadResults: misc.EmptyPodcastResults(false, nil),
	PodcastUrl:      "",

	LeftSide:    nil,
	PodcastList: nil,

	DebugState: nil,
	WhomBox:    nil,
	WhomButton: nil,

	TitleOrFnameRd: nil,

	TitleOrFnameBox:    nil,
	BackBox:            nil,
	StartDownloadBox:   nil,
	SelectAllOrNoneBox: nil,
	StopDownloadBox:    nil,

	ChosenTitles: nil,

	ChosenFnames: nil,
	KeyStream:    nil,
	SpinChar:     "|",
	Internetload: consts.DEFAULT_LOAD,

	Menue1Add:   nil,
	Menu2High:   nil,
	Menu2Medium: nil,
	Menu2Low:    nil,
	BusyCursor:  false,
}

func AllSelected(redrawWindow func(StateKind)) {
	if TheMediaWindow.ForceTitleOverFname {
		for index := range TheMediaWindow.ChosenTitles {
			TheMediaWindow.ChosenTitles[index] = true
		}
		redrawWindow(E_CHOOSEN_TITLES_MANY)
	} else {
		for index := range TheMediaWindow.ChosenFnames {
			TheMediaWindow.ChosenFnames[index] = true
		}
		redrawWindow(H_CHOOSEN_FILENAMES_MANY)
	}
}

func NoneSelected(redrawWindow func(StateKind)) {
	if TheMediaWindow.ForceTitleOverFname {
		for index := range TheMediaWindow.ChosenTitles {
			TheMediaWindow.ChosenTitles[index] = false
		}
		redrawWindow(C_CHOOSEN_TITLES_NONE)
	} else {
		for index := range TheMediaWindow.ChosenFnames {
			TheMediaWindow.ChosenFnames[index] = false
		}
		redrawWindow(F_CHOOSEN_FILENAMES_NONE)
	}
}

func SetRssType(dirName string) {
	rssFilePath := TheMediaWindow.ProgPath + "/" + dirName + "/" + consts.URL_OF_RSS_FN
	forceTitleOverFname, _, _ := podcasts.IsForceTitle(rssFilePath)
	TheMediaWindow.ForceTitleOverFname = forceTitleOverFname
	if forceTitleOverFname {
		TheMediaWindow.TitleOrFnameRd.SetSelected(TITLE_TYPE)
	} else {
		TheMediaWindow.TitleOrFnameRd.SetSelected(FILENAME_TYPE)
	}
}

const FILENAME_TYPE = "Use Filenames"
const TITLE_TYPE = "Use Episode Titles"

func FileOrTitle(redrawWindow func(StateKind)) *fyne.Container {
	var ftOptions = []string{FILENAME_TYPE, TITLE_TYPE}
	fileTitle := widget.NewRadioGroup(ftOptions, func(s string) {
		if s == FILENAME_TYPE {
			TheMediaWindow.ForceTitleOverFname = false
		} else {
			TheMediaWindow.ForceTitleOverFname = true
		}
		numberSelected := 0
		if TheMediaWindow.ForceTitleOverFname {
			for _, episodeSelected := range TheMediaWindow.ChosenTitles {
				if episodeSelected {
					numberSelected++
					if numberSelected > 1 {
						redrawWindow(E_CHOOSEN_TITLES_MANY)
						return
					}
				}
			}
		} else {
			for _, episodeSelected := range TheMediaWindow.ChosenFnames {
				if episodeSelected {
					numberSelected++
					if numberSelected > 1 {
						redrawWindow(H_CHOOSEN_FILENAMES_MANY)
						return
					}
				}
			}
		}

		if numberSelected == 1 {
			if TheMediaWindow.ForceTitleOverFname {
				redrawWindow(D_CHOOSEN_TITLES_ONE)
			} else {
				redrawWindow(G_CHOOSEN_FILENAMES_ONE)
			}
			return
		}
		if TheMediaWindow.ForceTitleOverFname {
			redrawWindow(C_CHOOSEN_TITLES_NONE)
		} else {
			redrawWindow(F_CHOOSEN_FILENAMES_NONE)
		}
	})
	fileTitle.Horizontal = true
	if TheMediaWindow.ForceTitleOverFname {
		fileTitle.Selected = TITLE_TYPE
	} else {
		fileTitle.Selected = FILENAME_TYPE
	}
	TheMediaWindow.TitleOrFnameRd = fileTitle
	horRadio := container.NewHBox(fileTitle)
	return horRadio
}

func ManyOneNone(redrawWindow func(StateKind)) {
	numberSelected := 0
	if TheMediaWindow.ForceTitleOverFname {
		for _, episodeSelected := range TheMediaWindow.ChosenTitles {
			if episodeSelected {
				numberSelected++
			}
		}
	} else {
		for _, episodeSelected := range TheMediaWindow.ChosenFnames {
			if episodeSelected {
				numberSelected++
			}
		}
	}

	if numberSelected == 0 {
		TheMediaWindow.StartDownloadBox.Hidden = true
		return
	}
	TheMediaWindow.StartDownloadBox.Hidden = false
}
