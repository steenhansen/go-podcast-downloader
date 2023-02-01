package terminal

import (
	"fmt"
	"testing"

	"github.com/steenhansen/go-podcast-downloader-console/src/consts"
	"github.com/steenhansen/go-podcast-downloader-console/src/misc"
	"github.com/steenhansen/go-podcast-downloader-console/src/terminal"
	"github.com/steenhansen/go-podcast-downloader-console/src/testings"
)

//  go run ./ siberiantimes.com/ecology/rss/ ecology

func TestAddByUrlAndName(t *testing.T) {

	url := consts.TEST_DIR_URL + "AddByUrlAndName/git-server-source/add-by-url-and-name.rss"

	progPath := misc.CurDir()
	testDir := progPath + "/Ecology"
	testings.DirRemove(testDir)

	podcastUrl := url //"siberiantimes.com/ecology/rss/"
	//	osArgs := []string{"test-prog.exe", "siberiantimes.com/ecology/rss/", "Ecology"}
	osArgs := []string{"AddByUrlAndName.exe", url, "Ecology"}
	progBounds := testings.ProgBounds(misc.CurDir())
	simKeyStream := make(chan string)

	podReport, err := terminal.AddByUrlAndName(podcastUrl, osArgs, progBounds, simKeyStream)
	fmt.Println("podReport", podReport, err)

	if podReport != "as;ldfh" {
		t.Fatalf(`TestAddByUrlAndName failed`)
	}
}
