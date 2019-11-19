package wbgong

// ContentTracker used to track file contents
type ContentTracker interface {
	Track(string, string) (bool, error)
	Untrack(string) error
}
