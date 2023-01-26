package rss

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

// no title is ok, if user gives us a title
func RssTitle(rss []byte) (string, error) { // 2 or 3
	type XmlTitle struct {
		Title string `xml:"channel>title"`
	}
	chanl := XmlTitle{Title: ""}
	err := xml.Unmarshal([]byte(rss), &chanl)
	if err != nil {
		return "", flaws.MissingTitle
	}
	title := strings.TrimSpace(chanl.Title)
	if len(title) == 0 {
		return "", flaws.EmptyTitle
	}
	return chanl.Title, nil
}

func RssItems(rss []byte) ([]string, []int, error) { // 2 or 3
	type XmlAttrib struct {
		UrlKey string `xml:"url,attr"`
		LenKey string `xml:"length,attr"`
	}
	type XmlEnclosures struct {
		Enclosures []XmlAttrib `xml:"channel>item>enclosure"`
	}
	enclosures := XmlEnclosures{}
	err := xml.Unmarshal([]byte(rss), &enclosures)
	if err != nil {
		return nil, nil, err
	}
	urls := make([]string, len(enclosures.Enclosures))
	lengths := make([]int, len(enclosures.Enclosures))
	if len(urls) == 0 {
		return nil, nil, flaws.EmptyItems
	}
	for i, v := range enclosures.Enclosures {
		urls[i] = string(v.UrlKey)
		length, err := strconv.Atoi(v.LenKey)
		if err != nil {
			lengths[i] = 0
		} else {
			lengths[i] = length
		}
	}
	return urls, lengths, nil
}
