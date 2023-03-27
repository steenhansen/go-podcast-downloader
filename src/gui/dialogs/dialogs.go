package dialogs

import (
	"image/color"
	"net/url"

	"podcast-downloader/src/gui/colors"
	"podcast-downloader/src/gui/values"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const INIT_WINDOW_WIDTH = 800
const INIT_WINDOW_HEIGHT = 400

const WINDOW_DEFAULT_TITLE = "Backup Podcasts"

func StartApp() fyne.Window {
	theApp := app.New()
	theApp.Settings().SetTheme(&myTheme{})
	fyneWindow := theApp.NewWindow(values.WINDOW_DEFAULT_TITLE)
	fyneWindow.Resize(fyne.NewSize(INIT_WINDOW_WIDTH, INIT_WINDOW_HEIGHT))
	TheMenu(fyneWindow)
	return fyneWindow
}

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	if name == theme.ColorNameDisabled {
		if variant == theme.VariantLight {
			return colors.DISABLED_LIGHT_CHECKBOX_TEXT
		}
		return colors.DISABLED_DARK_CHECKBOX_TEXT
	}

	if name == theme.ColorNameButton {
		if variant == theme.VariantLight {
			return colors.BUTTON_LIGHT_COLOR
		}
		return colors.BUTTON_DARK_COLOR
	}

	if name == theme.ColorNameForeground {
		if variant == theme.VariantLight {
			return colors.BLACK_TEXT
		}
		return colors.WHITE_TEXT
	}

	if name == theme.ColorNameDisabledButton {
		return theme.BackgroundColor()
	}

	if name == theme.ColorNamePrimary {
		return colors.SELECTED_RADIO_CHEK
	}

	if name == theme.ColorNameInputBorder {
		return colors.CHECKBOX_PERIMETER
	}

	if name == theme.ColorNameFocus {
		return colors.CURRENT_FOCUS_CIRCLE
	}

	if name == theme.ColorNameHover {
		return colors.HOVERING_CIRCLE
	}

	if name == theme.ColorNameInputBackground {
		return colors.UN_SELECTED_RADIO_CHEK
	}

	return theme.DefaultTheme().Color(name, variant)
}

func cancelAdd(podcastUrl string, fyneWindow fyne.Window, err error) {
	errorMess := err.Error()
	urlMess := podcastUrl + " is not a valid podcast"
	urlLabel := widget.NewLabel(urlMess)
	urlLabel.Alignment = fyne.TextAlignCenter
	cancelDialog := dialog.NewCustom(errorMess, "OK", urlLabel, fyneWindow)
	cancelDialog.Show()
}

func aboutDialog(fyneWindow fyne.Window) {
	aboutMess := "Podcast-Downloader"
	var gitRepo, _ = url.Parse("https://github.com/steenhansen/go-podcast-downloader")
	aboutLabel1 := widget.NewLabel("This is both a console and Windows")
	aboutLabel2 := widget.NewLabel("podcast downloader written in Go 1.20")
	hyperlink3 := widget.NewHyperlink("GitHub Repository", gitRepo)
	aboutLabels := container.NewVBox()
	aboutLabels.Add(aboutLabel1)
	aboutLabels.Add(aboutLabel2)
	aboutLabels.Add(hyperlink3)
	cancelDialog := dialog.NewCustom(aboutMess, "OK", aboutLabels, fyneWindow)
	cancelDialog.Show()
}
