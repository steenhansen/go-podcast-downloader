package initialize

import (
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/misc"
	"podcast-downloader/src/dos/podcasts"
)

const NASA_URL = "www.nasa.gov/rss/dyn/lg_image_of_the_day.rss"
const NASA_DIRNAME = "NASA Image of the Day"

func AddNasa() {
	progPath := misc.CurDir()
	podDirNames, _, _, _ := podcasts.AllPodcasts(progPath)
	if len(podDirNames) == 0 {
		forcingTitle := true
		media.ReSaveFolder(forcingTitle, progPath, NASA_DIRNAME, NASA_URL)
	}
}
