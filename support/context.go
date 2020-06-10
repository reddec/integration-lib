package support

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

var (
	rootContextInitializer sync.Once
	rootContext            context.Context
)

// Create global context that will be closed on SIGINT or SIGKILL signal from OS.
//
// The context will be created only once and it is safe to invoke it multiple time in a different go-routines
func SignalContext() context.Context {
	rootContextInitializer.Do(func() {
		ctx, closer := context.WithCancel(context.Background())
		go func() {
			c := make(chan os.Signal, 2)
			signal.Notify(c, os.Kill, os.Interrupt)
			for range c {
				closer()
				break
			}
		}()
		rootContext = ctx
	})
	return rootContext
}
