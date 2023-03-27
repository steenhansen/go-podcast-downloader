package selecting

import (
	"bufio"
	"fmt"
	"io"
	"podcast-downloader/src/dos/misc"

	"log"

	"os"

	"fyne.io/fyne/v2"
)

// https://developer.fyne.io/extend/bundle
func ButtonIcon(fileName string) *fyne.StaticResource {
	curDir := misc.CurDir()
	filePath := curDir + "/src/gui/images/" + fileName
	iconFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(iconFile)

	b, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Using dynamic icons, with ButtonIcon() and const UseDyanmicButtonIcons, not bundled icons, for ", fileName)
	return fyne.NewStaticResource("icon", b)
}

/*

 fyne bundle -o ./src/gui/selecting/res-prog-icon.go        --package selecting ./src/gui/images/prog-icon.png
 fyne bundle -o ./src/gui/selecting/res-go-back.go          --package selecting ./src/gui/images/go-back.png
 fyne bundle -o ./src/gui/selecting/res-select-all.go       --package selecting ./src/gui/images/select-all.png
 fyne bundle -o ./src/gui/selecting/res-select-none.go      --package selecting ./src/gui/images/select-none.png
 fyne bundle -o ./src/gui/selecting/res-stop-downloading.go --package selecting ./src/gui/images/stop-downloading.png

 fyne bundle -o ./src/gui/redux/res-busy-cursor.go --package redux ./src/gui/images/busy-cursor.png

*/
