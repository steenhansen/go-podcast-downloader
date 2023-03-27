package main

import (
	"fmt"
	"podcast-downloader/src/dos/globals"
	"podcast-downloader/src/gui/dialogs"
	"podcast-downloader/src/gui/redux"
	"podcast-downloader/src/gui/values"
)

func main() { // GUI version

	if values.UseDyanmicButtonIcons {
		fmt.Println("Debugging A - values.UseDyanmicButtonIcons = true ")
	}
	if values.GUI_DEBUG {
		fmt.Println("Debugging B - values.GUI_DEBUG = true ")
	}

	if globals.LogChannels {
		fmt.Println("Debugging C - globals.LogChannels = true ")
	}

	if globals.LogMemory {
		fmt.Println("Debugging D - globals.LogMemory = true ")
	}
	if globals.ForceTitle {
		fmt.Println("Debugging E - globals.ForceTitle = true ")
	}
	if globals.DnsErrorsTest {
		fmt.Println("Debugging F - globals.DnsErrorsTest = true ")
	}
	if globals.EmptyFilesTest {
		fmt.Println("Debugging G - globals.EmptyFilesTest = true ")
	}

	fyneWindow := dialogs.StartApp()

	redux.WindowStart(fyneWindow)
	fyneWindow.ShowAndRun()

}
