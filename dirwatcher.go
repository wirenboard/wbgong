package wbgo

import "time"

// DirWatcherClient is a consumer of FS changes
type DirWatcherClient interface {
	LoadFile(path string) error
	LiveLoadFile(path string) error
	LiveRemoveFile(path string) error
}

// DirWatcher provide FS changes interface
type DirWatcher interface {
	SetDelay(time.Duration)
	Load(string) error
	Stop()
}
