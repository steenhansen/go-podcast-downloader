package dialogs

import (
	"podcast-downloader/src/dos/consts"
	"podcast-downloader/src/gui/add"
	"podcast-downloader/src/gui/state"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func TheMenu(fyneWindow fyne.Window) {
	dialogTitle := "Add Podcast RSS Feed URL e.g. www.nasa.gov/rss/dyn/lg_image_of_the_day.rss                         "
	rssUrlEntry := widget.NewEntry()
	rssUrlEntry.Text = ""
	rssUrlItem := widget.NewFormItem("RSS URL", rssUrlEntry)
	formItems := []*widget.FormItem{rssUrlItem}
	addDialog := dialog.NewForm(dialogTitle, "OK", "Cancel", formItems,
		func(okChosen bool) {
			if okChosen {
				podcastUrl := strings.TrimSpace(rssUrlEntry.Text)
				if podcastUrl != "" {
					err := add.AddUrl(podcastUrl)
					if err != nil {
						cancelAdd(podcastUrl, fyneWindow, err)
					}
				}
			}
		}, fyneWindow)

	menu1Add := fyne.NewMenuItem("Add Podcast Url", func() { addDialog.Show() })
	menu2High := fyne.NewMenuItem("High Load - interferes with YouTube", nil)
	menu2Medium := fyne.NewMenuItem("Medium Load", nil)
	menu2Low := fyne.NewMenuItem("Low Load", nil)

	menu2High.Action = func() {
		state.TheMediaWindow.Internetload = consts.HIGH_LOAD
		menu2High.Checked = true
		menu2Medium.Checked = false
		menu2Low.Checked = false
	}

	menu2Medium.Action = func() {
		state.TheMediaWindow.Internetload = consts.MEDIUM_LOAD
		menu2High.Checked = false
		menu2Medium.Checked = true
		menu2Low.Checked = false
	}

	menu2Low.Action = func() {
		state.TheMediaWindow.Internetload = consts.LOW_LOAD
		menu2High.Checked = false
		menu2Medium.Checked = false
		menu2Low.Checked = true
	}

	menu3About := fyne.NewMenuItem("About", func() { aboutDialog(fyneWindow) })
	menu1 := fyne.NewMenu("File", menu1Add)
	menu2 := fyne.NewMenu("Internet Load", menu2High, menu2Medium, menu2Low)

	menu3 := fyne.NewMenu("About", menu3About)
	wholeMenu := fyne.NewMainMenu(menu1, menu2, menu3)
	fyneWindow.SetMainMenu(wholeMenu)

	state.TheMediaWindow.Menue1Add = menu1Add
	state.TheMediaWindow.Menu2High = menu2High
	state.TheMediaWindow.Menu2Medium = menu2Medium
	state.TheMediaWindow.Menu2Low = menu2Low

}
