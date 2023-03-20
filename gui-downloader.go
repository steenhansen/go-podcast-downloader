package main

import (
	"podcast-downloader/src/gui/dialogs"
	"podcast-downloader/src/gui/redux"
)

func main() {

	fyneWindow := dialogs.StartApp()
	redux.WindowStart(fyneWindow)
	fyneWindow.ShowAndRun()
}
