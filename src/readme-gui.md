

Build without console
  go build -ldflags -H=windowsgui .

Build with debugging
  go build .



https://developer.fyne.io/explore/redux


https://blogvali.com/class-5-fyne-golang-gui-course-checkbox/

https://developer.fyne.io/container/box







func red_button() *fyne.Container { 
    btn := widget.NewButton("Visit", nil)
    btn_color := canvas.NewRectangle(color.NRGBA{R: 0, G: 255, B: 0, A: 255})
    container1 := container.New(
        layout.NewMaxLayout(),
        btn_color,
        btn,
    )
    return container1
}