package finalizer

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"
)

// Cancel context that ensures that AfterFunc is called in event of ctrl-c or runtime error
//
// Add this to the beginning of main()
//		ctx, cancel := finalizer.WithCancel(context.Background())
//		defer func() { cancel(); <-ctx.Done() }()
// And this to the end
//		<-ctx.Done()

type finalizerContext struct {
	ctx                context.Context
	afterFuncWaitGroup sync.WaitGroup
}

func (f *finalizerContext) Deadline() (deadline time.Time, ok bool) {
	return f.ctx.Deadline()
}

func (f *finalizerContext) Done() <-chan struct{} {
	select {
	case <-f.ctx.Done():
		return f.ctx.Done()
	default:
		ch := make(chan struct{})
		go func() {
			<-f.ctx.Done()
			f.afterFuncWaitGroup.Wait()
			close(ch)
		}()
		return ch
	}
}

func (f *finalizerContext) Err() error {
	return f.ctx.Err()
}

func (f *finalizerContext) Value(key any) any {
	return f.ctx.Value(key)
}

func WithCancel(parentContext context.Context) (ctx *finalizerContext, cancel func()) {
	ctx = &finalizerContext{}
	ctx.ctx, cancel = context.WithCancel(parentContext)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		select {
		case <-sigChan:
			signal.Reset()
			cancel()
		case <-ctx.ctx.Done():
		}
	}()
	return
}

func AfterFunc(ctx *finalizerContext, f func()) {
	ctx.afterFuncWaitGroup.Add(1)
	context.AfterFunc(ctx.ctx, func() {
		f()
		ctx.afterFuncWaitGroup.Done()
	})
}
