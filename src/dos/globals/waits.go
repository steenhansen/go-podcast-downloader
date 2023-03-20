package globals

import (
	"sync"
)

type WaitCountDebug struct {
	currentCount int
}

var WaitCount = WaitCountDebug{currentCount: 0}

var waitMutex = sync.RWMutex{}

func (waitCount *WaitCountDebug) Adding() {
	waitMutex.Lock()
	waitCount.currentCount++
	waitMutex.Unlock()
}

func (waitCount *WaitCountDebug) Current() int {
	waitMutex.Lock()
	currentCount := waitCount.currentCount
	waitMutex.Unlock()
	return currentCount
}

func (waitCount *WaitCountDebug) Subtracting() {
	waitMutex.Lock()
	waitCount.currentCount--
	waitMutex.Unlock()
}
