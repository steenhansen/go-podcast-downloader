package consts

const LOG_CHANNELS = `^\-\-logChannels`

const FORCE_TITLE = `^\-\-forceTitle`
const MAX_TITLE_LEN = 80
const OPTION_FORCE_TITLE = `--forceTitle`

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

const MINIMUM_PLAIN = `^\-\-minimumDisk`
const MINIMUM_AND_NUMBER = `^\-\-minimumDisk=(\d|_)+`
const MINIMUM_NUMBER = `(\d|_)+`
