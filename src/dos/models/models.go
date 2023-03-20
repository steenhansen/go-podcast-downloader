package models

import (
	"context"
	"net/http"
	"time"
)

type MediaError struct {
	EnclosureUrl  string
	EnclosurePath string
	OrgErr        error
}

type MediaEnclosure struct {
	EnclosureUrl  string
	EnclosurePath string
	EnclosureSize int
}

type ProgBounds struct {
	ProgPath    string
	LoadOption  string
	LimitOption int
	MinDisk     int
	LogChannels bool
}

type CurStat struct {
	MinDiskMbs  int
	NetworkLoad string
}

type PodcastData struct {
	PodTitle  string
	PodPath   string
	PodUrls   []string
	PodSizes  []int
	PodTitles []string
}

type PodcastResults struct {
	ReadFiles     int
	SavedFiles    int
	PossibleFiles int
	VarietyFiles  string
	PodcastTime   time.Duration
	WasCanceled   bool
	SeriousError  error
}

type ReadLineFn func() string

type HttpFn func(ctx context.Context, mediaUrl string, numRetries int) (*http.Response, error)

// https://scene-si.org/2018/08/06/basic-monitoring-of-go-apps-with-the-runtime-package/
type MonitorMem struct {
	Current    uint64
	Cumulative uint64
	System     uint64
}
