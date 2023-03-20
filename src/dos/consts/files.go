package consts

import "strings"

const KB_BYTES int64 = 1024
const MB_BYTES int64 = 1024 * 1024
const GB_BYTES int64 = 1024 * 1024 * 1024
const TB_BYTES int64 = 1024 * 1024 * 1024 * 1024

const URL_OF_RSS_FN = "_origin-rss-url"
const SOURCE_FOLDER_TERMINAL = "dos"
const SOURCE_FOLDER_GUI = "gui"

const BAD_FILE_CHAR_AND_DOT = `[\\/:"*?<>|.]+`

func IsTesting(osArgs []string) bool {
	//	return false
	for _, anArg := range osArgs {
		if strings.HasPrefix(anArg, TEST_FLAG_PREFIX) {
			return true
		}
	}
	return false
}
