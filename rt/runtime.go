package rt

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/sunary/kitchen/l"
)

var (
	ll = l.Logger{}
)

// HandleCrash ...
func HandleCrash(handlers ...func(interface{})) {
	if r := recover(); r != nil {
		debug.PrintStack()
		logPanic(r)

		for _, fn := range handlers {
			fn(r)
		}
	}
}

func logPanic(r interface{}) {
	callers := ""
	for i := 0; true; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		callers = callers + fmt.Sprintf("%v:%v\n", file, line)
	}

	ll.Error("Recovered from panic", l.Object("e", r), l.String("caller", callers))
}

// RunWorker ...
func RunWorker(f func(), handlers ...func(interface{})) {
	go func() {
		defer HandleCrash(handlers...)

		f()
	}()
}

// RunWorkerUtilStop ...
func RunWorkerUtilStop(f func(), stopCh <-chan struct{}, handlers ...func(interface{})) {
	go func() {
		for {
			select {
			case <-stopCh:
				return

			default:
				func() {
					defer HandleCrash(handlers...)
					f()
				}()
			}
		}
	}()
}
