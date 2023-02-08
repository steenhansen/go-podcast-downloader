package models

import (
	"context"
	"net/http"
	"time"
)

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
}

type CurStat struct {
	ReadFiles   *int
	SavedFiles  *int
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
	Err           error
}

type ReadLineFn func() string

type HttpFn func(ctx context.Context, mediaUrl string) (*http.Response, error)
