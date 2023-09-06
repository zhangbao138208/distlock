package mutex

import (
	"context"
	"sync"
	"time"
)

type internal struct {
	locker   locker
	localMtx sync.Locker
}

func newInternal(locker locker, localMtx sync.Locker) internal {
	return internal{locker: locker, localMtx: localMtx}
}

func (itn *internal) Lock() {
	//itn.localMtx.Lock()
	itn.locker.lock()
	//itn.localMtx.Unlock()
}

func (itn *internal) LockCtx(ctx context.Context) bool {
	//itn.localMtx.Lock()
	//defer itn.localMtx.Unlock()
	return itn.locker.lockCtx(ctx)
}

func (itn *internal) TryLock() bool {
	//itn.localMtx.Lock()
	//defer itn.localMtx.Unlock()
	return itn.locker.tryLock()
}

func (itn *internal) Touch() bool {
	//itn.localMtx.Lock()
	//defer itn.localMtx.Unlock()
	return itn.locker.touch()
}

func (itn *internal) Unlock() {
	//itn.localMtx.Lock()
	itn.locker.unlock()
	//itn.localMtx.Unlock()
}

func (itn *internal) Until() time.Duration {
	//itn.localMtx.Lock()
	//defer itn.localMtx.Unlock()
	return time.Until(itn.locker.until)
}

func (itn *internal) Heartbeat(ctx context.Context) <-chan struct{} {
	notify := make(chan struct{})
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(itn.Until()):
				if !itn.Touch() {
					close(notify)
					return
				}
			}
		}
	}()
	return notify
}
