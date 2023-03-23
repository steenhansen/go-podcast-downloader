

# Fyne Gui & Console Podcast Downloader written in Go


## Fyne Gui Version
![Selecting NASA Image of the Day episodes](/src/gui/images/gui-nasa.png)

## Console Version
![Console Menu](/src/dos/images/menu.png)




## For Backing-Up Podcasts

  -  Joe Rogan's podcasts are now only on Spotify, except for the single #1109 in the RSS feed https://joeroganexp.libsyn.com/rss. They used to be all publically available.

  -  The MP3 files of http://feeds.feedburner.com/PhpRoundtable are no longer accessible.

  - The Hardcore History podcast http://feeds.feedburner.com/dancarlin/history, only lists the last 14 episodes out of 69 released.


## For Checking Podcast Files
      
  - Check the existance of episode files
  
  - Verify filenames and/or episode titles and extensions

  - Make sure RSS XML file valid of "If Books Could Kill"
~~~
    > go run ./console-downloader.go https://feeds.buzzsprout.com/2040953.rss
~~~

  - See if episode hit counters work

## Finding Podcast Feeds
  - Use <a href='https://castos.com/tools/find-podcast-rss-feed/'>castos.com</a> to find urls of feeds


## Adding a Podcast via the menu via "File | Add Podcast Url"

![Console Menu](/src/gui/images/add-rss.png)

  - American Scandal 
          
        rss.art19.com/american-scandal


  - BBC News Top stories
         
        podcasts.files.bbci.co.uk/p02nq0gn.rss?_BBC_News_Top_Stories_


  - BC Today from CBC Raido British Columbia

        www.cbc.ca/podcasting/includes/bcalmanac.xml

  - Black Box Down
      
        feeds.megaphone.fm/blackboxdown


 -  Breaking Points with Krystal and Saagar

         feeds.megaphone.fm/BRPL9803447123?_Breaking_Points_ 

  - English News - NHK WORLD RADIO JAPAN

         www3.nhk.or.jp/rj/podcast/rss/english.xml


  - Heist Podcast

         heistpodcast.libsyn.com/rss

  - Ecology (Siberian Times)

         siberiantimes.com/ecology/rss/

  - Nasa Image of the Day
      
        www.nasa.gov/rss/dyn/lg_image_of_the_day.rss 

  - Stuff You Should Know

         omnycontent.com/d/playlist/e73c998e-6e60-432f-8610-ae210140c5b1/A91018A4-EA4F-4130-BF55-AE270180C327/44710ECC-10BB-48D1-93C7-AE270180C33E/podcast.rss?_Stuff_You_Should_Know_

  - The Dollop with Dave Anthony and Gareth Reynolds
    
         www.omnycontent.com/d/playlist/885ace83-027a-47ad-ad67-aca7002f1df8/22b063ac-654d-428f-bd69-ae2400349cde/65ff0206-b585-4e2a-9872-ae240034c9c9/podcast.rss?_The_Dollop_


  - The History of the Twentieth Century
        
        history20th.libsyn.com/rss

  - The Rest Is History 

         rss.acast.com/the-rest-is-history-podcast

  - The SFFaudio Public Domain PDF Page (example of missing files)

         sffaudio.herokuapp.com/pdf/rss

  - Timesuck with Dan Cummins

        feeds.simplecast.com/Llc7KL2K?_Timesuck_with_Dan_Cummins_

  - Witness History

        podcasts.files.bbci.co.uk/p004t1hd.rss?_Witness_History_BBC_ 

  - Dan Carlin's Hardcore History

        http://feeds.feedburner.com/dancarlin/history

### [Compiling Windows Desktop Version](./src/gui.md)

### [Fyne GUI Resources](./src/resources.md)

### [Compiling Command Line Version](./src/dos.md)

### [Testing](./src/testing.md)









