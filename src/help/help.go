package help

func HelpText() string {
	help := `

This desktop program downloads podcast episodes onto your computer.
The program can be used to verify RSS feeds, or to collect files 
and ensure future use like with the old 'The Joe Rogan Experience'
RSS feed of http://joeroganexp.joerogan.libsynpro.com/rss

Run with menu :
  > go run ./
      1 |              jpg |  10 files |    1MB | Ecology
      2 |              mp3 |   1 files |    6MB | English News - NHK WORLD RADIO JAPAN
      3 |     jpg jpeg png |  40 files |  192MB | My Favorite Image Podcast
      4 |              mp3 |   6 files |  167MB | Reading Short and Deep
      5 |              mp3 |  11 files |  340MB | Stuff You Should Know
      6 |              pdf |  17 files |    8MB | The SFFaudio Public Domain PDF Page
      'Q' or a number + enter: 3

Add single podcast named "My Favorite Image Podcast" :
  > go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss My Favorite Image Podcast

Add single podcast named with podcast title :
  > go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss

Update single podcast
  > go run ./ My Favorite Image Podcast

Delete a podcast :
  > rmdir 'My Favorite Image Podcast'

Options :
  networkLoad (high | medium | low):
  > go run ./ --networkLoad=high

  fileLimit (integer):
  > go run ./ --fileLimit=3

  create empty files for filename checking
  > go run ./ --emptyFiles



--forceTitle

	`
	return help
}
