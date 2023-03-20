package rss

// no title is ok, if user gives us a title
type xmlRssTitle struct {
	Title string `xml:"channel>title"`
}

type xmlItemTitles struct {
	Titles []string `xml:"channel>item>title"` // itunes:title are ignored by changing them to itunes:tiXYZ
}

type xmlUrlLen struct {
	UrlKey string `xml:"url,attr"`
	LenKey string `xml:"length,attr"`
}

type xmlEnclosures struct {
	Enclosures []xmlUrlLen `xml:"channel>item>enclosure"`
}
