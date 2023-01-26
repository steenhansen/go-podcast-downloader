package help

func HelpText() string {
	help := `



go run ./
  - gives user numbered rss feed choices

go run ./ help
  - shows help

[IF DOES NOT EXIST]
go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
  - adds 'NASA Image of the Day' to podcasts

[IF DOES NOT EXIST]
go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss todays nasa image
  - adds 'todays nasa image' to podcasts

[IF EXISTS]
go run ./ https://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss
  - downloads new images

[IF EXISTS]
go run ./ todays nasa image
  - downloads new images


	`
	return help
}
