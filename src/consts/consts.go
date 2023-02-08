package consts

import (
	"time"
)

const URL_OF_RSS_FN = "_origin-rss-url"
const SOURCE_FOLDER = "src"

const MAX_READ_FILE_TIME = time.Millisecond * 400000000

const KB_BYTES int64 = 1024
const MB_BYTES int64 = 1024 * 1024
const GB_BYTES int64 = 1024 * 1024 * 1024
const TB_BYTES int64 = 1024 * 1024 * 1024 * 1024

const MIN_DISK_BYTES int = 1_000_000_000
const MIN_DISK_FAIL_BYTES int = 999_000_000_000_000

const BAD_FILE_CHAR_AND_DOT = `[\\/:"*?<>|.]+`

const EMTPY_FILES = `^\-\-emptyFiles`

const LIMIT_PLAIN = `^\-\-fileLimit`
const LIMIT_AND_NUMBER = `^\-\-fileLimit=\d+`
const LIMIT_NUMBER = `\d+`

const LOAD_PLAIN = `^\-\-networkLoad`
const LOAD_AND_SPEED = `^\-\-networkLoad=(high|medium|low)`
const LOAD_CHOICE = `high|medium|low`

const HIGH_LOAD = "high"
const MEDIUM_LOAD = "medium"
const LOW_LOAD = "low"

const HIGH_SLEEP = 0
const MEDIUM_SLEEP = 5
const LOW_SLEEP = 10
const CPUS_MED_DIVIDER = 4

const HELP_PLAIN = "help"
const HELP_DASH = "-help"
const HELP_DASH_DASH = "--help"

const CLEAR_SCREEN = "\033[H\033[2J"

const SINGLE_DASH_ALPHA = `$-\w*`
const RACE_DEBUG = "$-race"

const ERROR_PREFIX = "*** "

const STOP_KEY_LOWER = "s"
const QUIT_KEY_LOWER = "q"

const HTML_404_BEGIN = "<!DOCTYPE"

const TEST_FLAG_PREFIX = "-test."

const TEST_DIR_URL = "https://raw.githubusercontent.com/steenhansen/pod-down-consol/main/src/internet-tests/"

const KEY_BUFF_SIZE = 1
const KEY_BUFF_ERROR = "GetKeys() keyboard error"

const HTTP_OK_RESP = 200

const FIRST_BYTES_OF_ERROR_PAGE = 100
