package clockify

import (
	"time"
)

type Fs struct {
	API *API
}

func New(config *Config) *Fs {
	return &Fs{
		API: &API{},
	}
}

func (c Fs) ProjectFilepath(key string) string {
	panic("implement me")
}

func (c Fs) ProjectFilepaths() ([]string, error) {
	panic("implement me")
}

func (c Fs) RecordFilepath(start time.Time) string {
	panic("implement me")
}

func (c Fs) RecordFilepaths(dir string, less func(a string, b string) bool) ([]string, error) {
	panic("implement me")
}

func (c Fs) RecordDirs() ([]string, error) {
	panic("implement me")
}

func (c Fs) RecordDirFromDate(date time.Time) string {
	panic("implement me")
}

func (c Fs) EnsureDirectories() error {
	return nil
}

func (c Fs) EnsureRecordDir(date time.Time) error {
	return nil
}
