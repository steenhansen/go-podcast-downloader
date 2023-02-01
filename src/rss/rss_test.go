package rss

import (
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/flaws"
)

//      go test ./...

func TestRssNoTitle(t *testing.T) {
	none := []byte("no title")
	_, err := RssTitle(none)
	if err != flaws.MissingTitle {
		t.Fatalf(`TestRssNoTitle failed`)
	}
}

// func TestRssEmptyTitle(t *testing.T) {
// 	empty := []byte("<channel><title>   </title></channel>")
// 	_, err := RssTitle(empty)
// 	if err != flaws.EmptyTitle {
// 		t.Fatalf(`TestRssEmptyTitle failed`)
// 	}
// }

// func TestRssEmptyItems(t *testing.T) {
// 	//	none := []byte("<channel><title>a-title</title><item><enclosure></enclosure></item></channel>")
// 	none := []byte(
// 		`<?xml version="1.0" encoding="UTF-8"?>
// 	 <rss version="2.0">
// 	<channel>
// 	<title>NASA Image of the Day</title>
//   <item></item>
// </channel>
// </rss>
// `)
// 	_, _, err := RssItems(none)
// 	//	fmt.Println(" TestRssEmptyItems >>", res, err, EmptyItems)
// 	if err != flaws.EmptyItems {
// 		t.Fatalf(`TestRssEmptyItems failed`)
// 	}
// }

// func TestRssEmptyItems(t *testing.T) {
// 	//	none := []byte("<channel><title>a-title</title><item><enclosure></enclosure></item></channel>")
// 	none := []byte(
// 		`<?xml version="1.0" encoding="UTF-8"?>
// 	 <rss version="2.0">
// 	<channel> <title>NASA Image of the Day</title>
//  <description>The latest NASA &quot;Image of the Day&quot; image.</description>
//  <link>http://www.nasa.gov/</link>
//  <atom:link rel="self" href="http://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss" />
//  <language>en-us</language>
//  <managingEditor>brian.dunbar@nasa.gov</managingEditor>
//  <webMaster>brian.dunbar@nasa.gov</webMaster>
//  <docs>http://blogs.law.harvard.edu/tech/rss</docs>

//  <item> <title>Sun Rings in New Month with Strong Flare</title>
//  <link>http://www.nasa.gov/image-feature/sun-rings-in-new-month-with-strong-flare</link>
//  <description>The Sun released an X1 solar flare, captured by our Solar Dynamics Observatory (SDO) on Oct. 2, 2022.</description>
//  <enclosure url="http://www.nasa.gov/sites/default/files/thumbnails/image/oct-2-2022-x1-flare-131-171-1024x1024.jpeg" length="642777" type="image/jpeg" />
//  <guid isPermaLink="false">http://www.nasa.gov/image-feature/sun-rings-in-new-month-with-strong-flare</guid>
//  <pubDate>Tue, 11 Oct 2022 11:39 EDT</pubDate>
//  <source url="http://www.nasa.gov/rss/dyn/lg_image_of_the_day.rss">NASA Image of the Day</source>
// </item>
// </channel>
// </rss>
// `)
// 	res, err := rssItems(none)
// 	fmt.Println(" TestRssEmptyItems >>", res, err, EmptyItems)
// 	if err != EmptyItems {
// 		t.Fatalf(`TestRssEmptyItems failed`)
// 	}
// }

// func TestInitFolder(t *testing.T) {
// 	path := "x:/does-not-exist"
// 	title := "a-title"
// 	expect := flaws.CantCreateDirSerious.ContinueError(path+"/"+title, nil)
// 	_, _, err := media.InitFolder(path, title, "http://www.pod.cast")
// 	if err.Error() != expect.Error() {
// 		t.Fatalf(`TestInitFolder failed`)
// 	}
// }
