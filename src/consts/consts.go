package consts

import (
	"time"
)

const HTTP_RETRIES = 16 // go run ./ feeds.megaphone.fm/blackboxdown --forceTitle "had 4 episodes retry 8 times"

const MEDIA_MAX_READ_FILE_TIME = time.Hour * 1 // max wait per media episode file

const RSS_MAX_READ_FILE_TIME = time.Minute * 2 // max wait for single RSS XML file

const DEFAULT_LOAD = LOW_LOAD
