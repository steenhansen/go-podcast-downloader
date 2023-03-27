
# Running and Compiling GUI Version

## Requirements

  - Go 1.20
  - Fyne v2

## Only One main() Function Out of Two
```
gui-downloader.go
    func main() {  // GUI Version On
```
```
console-downloader.go
    func main___not_main() {  // Console Version Off
```
## Download Dependencies
```
    > go mod tidy
```

## Running GUI From Source

#### This takes a long time the first time
```
    > go run ./gui-downloader.go
```
## Compiling GUI Executable With Debug Console, No Icon
```
    > go build ./gui-downloader.go
```

## Compiling GUI Executable Without Console, No Icon
```
    > go build -ldflags -H=windowsgui ./gui-downloader.go
```

## Packaging Independent GUI Executable Without Dependencies, With Icon
```
    > fyne package -os windows -icon ./src/gui/images/prog-icon.png
```