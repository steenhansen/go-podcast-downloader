
# Running and Compiling Console Version

## Requirements

  - Go 1.20
  - Fyne v2

## Only One main() Function Out of Two
```
gui-downloader.go
    func main___not_main() {  // GUI Version Off
```
```
console-downloader.go
    func main() {  // Console Version On
```
## Download Dependencies
```
    > go mod tidy
```

## Running in Console From Source 
```
    > go run ./console-downloader.go
```
## Compiling Independent Console Executable, No Icon
```
    > go build ./console-downloader.go
```

