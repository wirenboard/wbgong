package wbgong

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const (
	deferredCapacity = 256
	kiB              = 1024
	stackBufferSize  = 32 * kiB
)

func doVisit(visitor, thing any, methodName string, args []any) bool {
	method, found := reflect.TypeOf(visitor).MethodByName(methodName)
	if !found {
		return false
	}

	moreValues := make([]reflect.Value, len(args))
	for i, arg := range args {
		moreValues[i] = reflect.ValueOf(arg)
	}
	method.Func.Call(append([]reflect.Value{
		reflect.ValueOf(visitor),
		reflect.ValueOf(thing),
	}, moreValues...))

	return true
}

func Visit(visitor, thing any, prefix string, args ...any) {
	typeName := reflect.Indirect(reflect.ValueOf(thing)).Type().Name()
	methodName := prefix + typeName
	if !doVisit(visitor, thing, methodName, args) &&
		!doVisit(visitor, thing, prefix+"Anything", args) {
		Debug.Printf("visit: no visitor method for %s", typeName)
		return
	}
}

type DeferredList struct {
	sync.Mutex
	fns      []func()
	executor func(func())
}

func NewDeferredList(executor func(func())) *DeferredList {
	return &DeferredList{
		fns:      make([]func(), 0, deferredCapacity),
		executor: executor,
	}
}

func (dl *DeferredList) MaybeDefer(thunk func()) {
	dl.Lock()
	if dl.fns != nil {
		dl.fns = append(dl.fns, thunk)
		dl.Unlock()
		return
	}
	dl.Unlock()
	if dl.executor != nil {
		dl.executor(thunk)
	} else {
		thunk()
	}
}

func (dl *DeferredList) Ready() {
	dl.Lock()
	fns := dl.fns
	dl.fns = nil
	dl.Unlock()
	for _, fn := range fns {
		fn()
	}
}

// Truename returns the shortest absolute pathname
// leading to the specified existing file.
// Note that a single file may have multiple
// hard links or be accessible via multiple bind mounts.
// Truename doesn't account for these situations.
func Truename(filePath string) (string, error) {
	p, err := filepath.EvalSymlinks(filePath)
	if err != nil {
		return filePath, fmt.Errorf("failed to evaluate symlinks: %w", err)
	}
	p, err = filepath.Abs(p)
	if err != nil {
		return filePath, fmt.Errorf("failed to get absolute path: %w", err)
	}
	return filepath.Clean(p), nil
}

// IsSubpath returns true if maybeSubpath is a subpath of basepath. It
// uses filepath to be compatible with os-dependent paths.
func IsSubpath(basepath, maybeSubpath string) bool {
	rel, err := filepath.Rel(basepath, maybeSubpath)
	if err != nil {
		return false
	}
	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		// not a subpath if "../" has to be used for the relative path
		return false
	}
	return true
}

func GetStack() string {
	buf := make([]byte, stackBufferSize)
	n := runtime.Stack(buf, true)
	return string(buf[0:n])
}
