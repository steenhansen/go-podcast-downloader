package add

import (
	"podcast-downloader/src/dos/media"
	"podcast-downloader/src/dos/podcasts"
	"podcast-downloader/src/dos/rss"
	"podcast-downloader/src/gui/redux"
	"podcast-downloader/src/gui/state"
)

func AddUrl(podcastUrl string) error {
	rssXml, _, _, _, err := podcasts.ReadRssUrl(podcastUrl, rss.HttpReal)
	if err != nil {
		redux.RedrawWindow(state.B_NOTHING_SELECTED)
		return err
	}
	mediaTitle, _ := rss.RssTitle(rssXml)
	forcingTitle := false
	progPath := state.TheMediaWindow.ProgPath
	media.ReSaveFolder(forcingTitle, progPath, mediaTitle, podcastUrl)
	redux.RedrawWindow(state.A_LOAD_DIRECTORIES)
	redux.RedrawWindow(state.B_NOTHING_SELECTED)
	return nil
}
