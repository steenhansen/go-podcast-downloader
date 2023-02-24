
# Can be used to "back-up" a podcast

  -  Joe Rogan's podcasts are only on Spotify, except for the single #1109 in the RSS feed https://joeroganexp.libsyn.com/rss.

  -  The MP3 files of http://feeds.feedburner.com/PhpRoundtable are no longer accessible.

  - The Hardcore History podcast http://feeds.feedburner.com/dancarlin/history, only lists the last 14 episodes out of 69.


aaaaaaaaaaa

![Roman 38 is square root of 1444](src/images/menu.png)

ddddddddddddd

![Roman 38 is square root of 1444](src/images/nasa-images.png)

ggggggggggggggggg

![Roman 38 is square root of 1444](src/images/nasa-rss.png)


This Go console program downloads podcast episodes to your local machine. 
Only new episodes are downloaded. Files are not erased. 


Can be used to test an RSS feed for missing files. 




To run with a menu
```
    go run  ./
```



Or to just download new episodes of one podcast


  - American Scandal
 
        go run ./ rss.art19.com/american-scandal
        go run ./ American Scandal


  - BC Today from CBC Raido British Columbia

        go run ./ www.cbc.ca/podcasting/includes/bcalmanac.xml
        go run ./ BC Today from CBC Raido British Columbia


  - BBC News Top stories

        go run ./ podcasts.files.bbci.co.uk/p02nq0gn.rss?_BBC_News_Top_Stories_
        go run ./ BBC News Top stories


  - BC Today from CBC Raido British Columbia

        go run ./ www.cbc.ca/podcasting/includes/bcalmanac.xml
        go run ./ BC Today from CBC Raido British Columbia


  - Black Box Down
    
        go run ./ feeds.megaphone.fm/blackboxdown --forceTitle
        go run ./ Black Box Down


 -  Breaking Points with Krystal and Saagar
     
        go run ./ feeds.megaphone.fm/BRPL9803447123?_Breaking_Points_ --forceTitle
        go run ./ Breaking Points with Krystal and Saagar


  - English News - NHK WORLD RADIO JAPAN

        go run ./ www3.nhk.or.jp/rj/podcast/rss/english.xml
        go run ./ English News - NHK WORLD RADIO JAPAN


  - Heist Podcast

        go run ./ heistpodcast.libsyn.com/rss
        go run ./ Heist Podcast


  - Nasa Image of the Day
  
        go run ./ www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
        go run ./ Nasa Image of the Day


  - Ecology (Siberian Times)

        > podcast-downloader.exe siberiantimes.com/ecology/rss/
        > podcast-downloader.exe Ecology


  - Stuff You Should Know
      
        go run ./ omnycontent.com/d/playlist/e73c998e-6e60-432f-8610-ae210140c5b1/A91018A4-EA4F-4130-BF55-AE270180C327/44710ECC-10BB-48D1-93C7-AE270180C33E/podcast.rss?_Stuff_You_Should_Know_
        go run ./ Stuff You Should Know


  - The Dollop with Dave Anthony and Gareth Reynolds

        go run ./ www.omnycontent.com/d/playlist/885ace83-027a-47ad-ad67-aca7002f1df8/22b063ac-654d-428f-bd69-ae2400349cde/65ff0206-b585-4e2a-9872-ae240034c9c9/podcast.rss?_The_Dollop_
        go run ./ The Dollop with Dave Anthony and Gareth Reynolds


  - The History of the Twentieth Century

        go run ./ history20th.libsyn.com/rss
        go run ./ The History of the Twentieth Century


  - The Rest Is History 
      
        go run ./ rss.acast.com/the-rest-is-history-podcast
        go run ./ The Rest Is History


  - The SFFaudio Public Domain PDF Page (example of missing files)
        
        go run ./ sffaudio.herokuapp.com/pdf/rss
        go run ./ The SFFaudio Public Domain PDF Page


  - Timesuck with Dan Cummins

        go run ./ feeds.simplecast.com/Llc7KL2K?_Timesuck_with_Dan_Cummins_
        go run ./ Timesuck with Dan Cummins


  - Witness History
    
        go run ./ podcasts.files.bbci.co.uk/p004t1hd.rss?_Witness_History_BBC_ 
        go run ./ Witness History





////////////////////////////////////////
////////////////////////////////////////
////////////////////////////////////////

FIND PODCAST FEEDS
  - Use castos.com to <a href='https://castos.com/tools/find-podcast-rss-feed/'>find urls of feeds</a>

TESTS
    go clean -testcache 

  - Test everything
```
    go test ./src/tests_mocked_http/... ./src/tests_real_internet/... -count=1
```
  - Test only real Internet tests
```
    go test ./src/tests_real_internet/... -count=1 -timeout 50s     

go test ./src/tests_real_internet/missing-file/... -count=1 -timeout 5s     OK


```
  - Test only mocked Internet tests
```
     go test ./src/tests_mocked_http/... -count=1
```


OPTIONAL ARGUMENTS
  - --forceTitle uses the title of each episode as the locally saved filename

      go run ./ --forceTitle

  - --networkLoad sets the amount of network traffic, default is "high"
      
      go run ./ --networkLoad=high
      
      go run ./ --networkLoad=medium
      
      go run ./ --networkLoad=low 

  - --fileLimit sets the maximum number of files to download from a podcast
    
      go run ./ --fileLimit=3

TESTING ARGUMENTS      

  - --emptyFiles podcasts are not actually read, used to check for file existance and filenames, all files are empty
    
      go run ./ --emptyFiles

  - --logChannels will save all channel signalling in /src/channelLog.txt for debugging
    
      go run ./ --logChannels


COMPILING EXECUTABLE
    
    > go build -o podcast-downloader.exe

    > ./podcast-downloader.exe


 
![Roman 38 is square root of 1444](src/images/channels.png)







