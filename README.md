


















RUN console version
  go run ./console-downloader.go




Run gui version
  go run ./gui-downloader.go

  
  
Compile gui version  
  go build ./gui-downloader.go



fyne package -os windows


https://developer.fyne.io/extend/bundle



 fyne bundle    -o ./src/gui/selecting/res-prog-icon.go           --package selecting ./src/gui/images/prog-icon.png
 fyne bundle    -o ./src/gui/selecting/res-go-back.go             --package selecting ./src/gui/images/go-back.png
 fyne bundle    -o ./src/gui/selecting/res-select-all.go          --package selecting ./src/gui/images/select-all.png
 fyne bundle    -o ./src/gui/selecting/res-select-none.go         --package selecting ./src/gui/images/select-none.png
 fyne bundle    -o ./src/gui/selecting/res-stop-downloading.go    --package selecting ./src/gui/images/stop-downloading.png



/////////////////////////////////////////////

https://developer.fyne.io/started/packaging


fyne package -os windows -icon ./src/gui/images/prog-icon.png


