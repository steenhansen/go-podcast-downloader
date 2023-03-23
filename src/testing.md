
# Testing

    go clean -testcache 

  - Test Terminal
    everything
```
    > go test ./src/dos/tests_mocked_http/... ./src/dos/tests_real_internet/... -count=1
```
  - Test only real Internet tests
```
    > go test ./src/dos/tests_real_internet/... -count=1 
```
  - Test only mocked Internet tests
```
    > go test ./src/dos/tests_mocked_http/... -count=1 
```

  - Test single test
```
    > go test ./src/dos/tests_real_internet/pressStop_r/... -count=1
```

## Optional Arguments
  - --forceTitle uses the title of each episode as the locally saved filename
  instead of the filename of the downloaded file which can be like Stitchers' "default.mp3".
  This is needed for "Black Box Down", "Breaking Points", and "Nasa Image of the Day".
```
    >  go run ./console-downloader.go --forceTitle
```

  - --networkLoad sets the amount of network traffic, default is "medium"
```      
    > go run ./console-downloader.go --networkLoad=high
      
    > go run ./console-downloader.go --networkLoad=medium
      
    > go run ./console-downloader.go --networkLoad=low 
```
  - --fileLimit sets the maximum number of files to download from a podcast
```
    > go run ./console-downloader.go --fileLimit=3
```
## Testing Arguments

  - --emptyFiles podcasts are not actually read, used to check for file existance and filenames only, all files are empty
```    
    > go run ./console-downloader.go --emptyFiles
```
  - --dnsErrors will ensure every http episode Get will have a DNS error first for debugging
```
    > go run ./console-downloader.go --dnsErrors
```

  - --minimumDisk will help test for running out of diskspace errors
```
    > go run ./console-downloader.go --minimumDisk=1_000_000_000_000_000
```

## Testing Global Variables
  - LogChannels==true will save all channel signalling in /src/channel-mem-log.txt for debugging

  - LogMemory==true will save memory use every minute in /src/channel-mem-log.txt for debugging


## Channels State Diagram
 
![How go routines, waitGroups, and channels interact](/src/dos/images/channels.png)







