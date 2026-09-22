package wbgong

import (
	"io"
	"log"
	"os"
	"sync"
)

var funcWriteFileAtomic = sync.OnceValue(func() func(string, io.Reader, os.FileMode) error {
	funcSym, errSym := plug.Lookup("WriteFileAtomic")
	if errSym != nil {
		log.Fatalf("Error in lookup symbol: %v", errSym)
	}
	fn, okResolve := funcSym.(func(string, io.Reader, os.FileMode) error)
	if !okResolve {
		log.Fatal("Wrong sign on resolving func")
	}
	return fn
})

// WriteFileAtomic leaves the old file intact until the complete replacement
// has been written, synced and closed. Existing readers keep the old inode.
// Symlinks are followed, and existing ownership and permissions are preserved.
// New files use perm, subject to umask. Directory syncing after replacement is
// best-effort; failures are logged rather than returned after the write commits.
func WriteFileAtomic(path string, content io.Reader, perm os.FileMode) error {
	return funcWriteFileAtomic()(path, content, perm)
}
