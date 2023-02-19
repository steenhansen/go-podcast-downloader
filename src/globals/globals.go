package globals

// Episode filenames will be forced to match <title> instead of actual downloaded file
// go run ./  feeds.megaphone.fm/BRPL9803447123?_Breaking_Points_ --forceTitle
var ForceTitle = false

// User pressed 's' to stop downloading, a variable instead of an error
var StopingOnSKey = false

// For testing, every http episode will first get a DnsError
// go run ./ feeds.megaphone.fm/blackboxdown --dnsErrors
var DnsErrorsTest = false

// For testing, forgo actually downloading episodes, just create a 0 byte file
// go run ./ rss.acast.com/the-rest-is-history-podcast --emptyFiles
var EmptyFilesTest = false
