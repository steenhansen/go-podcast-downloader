



go clean -testcache 

TESTS
  - Test everything
```
    go test ./... -count=1 | grep -E "^[^?]"
```
  - Tests using real internet 
```
    go test ./src/internet-tests/... -count=1 | grep -E "^[^?]"
```
  - Test using mocked internet
```
     go test ./... -run "^Test_.*$" -count=1 | grep -E "^[^?]"
```


OPTIONAL ARGUMENTS
  - --networkLoad sets the amount of network traffic
      
      go run ./ --networkLoad=high
      
      go run ./ --networkLoad=medium
      
      go run ./ --networkLoad=low 

  - --fileLimit sets the maximum number of files to download from a podcast
    
      go run ./ --fileLimit=3

  - --emptyFiles podcasts are not actually read, used to check for file existance, all files are empty
    
      go run ./ --emptyFiles


FIND PODCAST FEEDS
  - Use castos.com to <a href='https://castos.com/tools/find-podcast-rss-feed/'>find feeds</a>

TESTING WITH EMPTY FILES
  -  Stuff You Should Know

     go run ./ --emptyFiles https://omnycontent.com/d/playlist/e73c998e-6e60-432f-8610-ae210140c5b1/A91018A4-EA4F-4130-BF55-AE270180C327/44710ECC-10BB-48D1-93C7-AE270180C33E/podcast.rss 

  -  Nasa Image of the Day
  
      go run ./ --emptyFiles https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss

  -  The Rest Is History

      go run ./ --emptyFiles https://rss.acast.com/the-rest-is-history-podcast 

  -  Witness History - BBC

      go run ./ --emptyFiles https://podcasts.files.bbci.co.uk/p004t1hd.rss

  -  The SFFaudio Public Domain PDF Page (example of missing files)

      go run ./ --emptyFiles sffaudio.herokuapp.com/pdf/rss 

  -  The History of the Twentieth Century

      go run ./ --emptyFiles https://history20th.libsyn.com/rss

  -  The Dollop with Dave Anthony and Gareth Reynolds

      go run ./ --emptyFiles https://www.omnycontent.com/d/playlist/885ace83-027a-47ad-ad67-aca7002f1df8/22b063ac-654d-428f-bd69-ae2400349cde/65ff0206-b585-4e2a-9872-ae240034c9c9/podcast.rss

  -  English News - NHK WORLD RADIO JAPAN

      go run ./ --emptyFiles https://www3.nhk.or.jp/rj/podcast/rss/english.xml


