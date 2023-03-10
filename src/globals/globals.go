package globals

import "github.com/steenhansen/go-podcast-downloader/src/consts"

// Episode filenames will be forced to match <title> instead of actual downloaded file
// go run ./  feeds.megaphone.fm/BRPL9803447123?_Breaking_Points_ --forceTitle
var ForceTitle = false

// For testing, every http episode will first get a DnsError
// go run ./ feeds.megaphone.fm/blackboxdown --dnsErrors
var DnsErrorsTest = false

// For testing, forgo actually downloading episodes, just create a 0 byte file
// go run ./ rss.acast.com/the-rest-is-history-podcast --emptyFiles
var EmptyFilesTest = false

// For testing, logging all channel changes
// go run ./ rss.acast.com/the-rest-is-history-podcast -logChannels
var LogChannels = false

var MediaMaxReadFileTime = consts.MEDIA_MAX_READ_FILE_TIME
var RssMaxReadFileTime = consts.RSS_MAX_READ_FILE_TIME

// Default progBounds.MinDisk
var MinDiskBytes = consts.MIN_DISK_BYTES
