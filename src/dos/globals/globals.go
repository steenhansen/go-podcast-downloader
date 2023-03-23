package globals

import "podcast-downloader/src/dos/consts"

// Episode filenames will be forced to match <title> instead of actual downloaded file
// go run ./  feeds.megaphone.fm/BRPL9803447123?_Breaking_Points_ --forceTitle
var ForceTitle = false

var MediaMaxReadFileTime = consts.MEDIA_MAX_READ_FILE_TIME
var RssMaxReadFileTime = consts.RSS_MAX_READ_FILE_TIME

// For testing, every http episode will first get a DnsError
// go run ./ feeds.megaphone.fm/blackboxdown --dnsErrors
var DnsErrorsTest = false

// For testing, forgo actually downloading episodes, just create a 0 byte file
// go run ./ rss.acast.com/the-rest-is-history-podcast --emptyFiles
var EmptyFilesTest = false

// Default progBounds.MinDisk
var MinDiskBytes = consts.MIN_DISK_BYTES

// For testing, log all channel changes in /src/channel-mem-log.txt
var LogChannels = false

// For testing, log memory usage by the minute to /src/channel-mem-log.txt
var LogMemory = false
