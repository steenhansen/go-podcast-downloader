
# Resources

## Buttons as png files


### To debug with buttons as png files set below constant in /src/gui/values/values.go
```
const UseDyanmicButtonIcons = true
```



## Build Icons Into .Exe Resources
```
 fyne bundle -o ./src/gui/selecting/res-prog-icon.go        --package selecting ./src/gui/images/prog-icon.png
 fyne bundle -o ./src/gui/selecting/res-go-back.go          --package selecting ./src/gui/images/go-back.png
 fyne bundle -o ./src/gui/selecting/res-select-all.go       --package selecting ./src/gui/images/select-all.png
 fyne bundle -o ./src/gui/selecting/res-select-none.go      --package selecting ./src/gui/images/select-none.png
 fyne bundle -o ./src/gui/selecting/res-stop-downloading.go --package selecting ./src/gui/images/stop-downloading.png

 fyne bundle -o ./src/gui/redux/res-busy-cursor.go --package redux ./src/gui/images/busy-cursor.png


```


