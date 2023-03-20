package consts

import (
	"time"
)

const MEDIA_MAX_READ_FILE_TIME = time.Hour * 1 // max wait per media episode file

const RSS_MAX_READ_FILE_TIME = time.Minute * 2 // max wait for single RSS XML file

const DEFAULT_LOAD = HIGH_LOAD // default network load is --networkLoad=high

const MIN_DISK_BYTES int = 1_000_000_000 // 1 gb

const GOOD_BYE_MESS = "good bye"

const CHANNEL_LOG_NAME = "channelLog.txt"

const RETRIES_RSS_FILE = 2

const RETRIES_MEDIA_FILES = 16 // go run ./ feeds.megaphone.fm/blackboxdown --forceTitle "had 4 episodes retry 8 times"

const ASCII_CARRIAGE_RETURN = 13
