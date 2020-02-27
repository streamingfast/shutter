package shutter

import (
	"errors"
	"sync"
)

type Shutter struct {
	lock     sync.Mutex // shutdown lock

	// the error is assigned before any channel (terminating & terminated)
	// are closed. thus when calling `Err()` is *always* available.
	err                 error
	once                sync.Once

	// terminating occurs when the Shutdown() function is called on the shutter. Is is an opportunity to clean up loose ends
	// the terminating channel, is a signal to know that the process is shutting down
	terminatingCh chan struct{}
	terminatingFunc     []func(error)
	terminatingFuncLock sync.Mutex

	// terminated occurs when the Shutdown() function has been called and all the clean-up functions are completed. We can assume that
	// we can kill the process now without any negative effects.
	terminatedCh chan struct{}
	terminatedFunc     []func(error)
	terminatedFuncLock sync.Mutex
}

func New(opts ...Option) *Shutter {
	s := &Shutter{
		terminatingCh: make(chan struct{}),
		terminatedCh:make(chan struct{}),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Option = func(shutter *Shutter)


// registers a function to be called on terminated
func RegisterOnTerminated(f func(error)) Option {
	return func(s *Shutter) {
		s.OnTerminated(f)
	}
}

// registers a function to be called on terminating
func RegisterOnTerminating(f func(error)) Option {
	return func(s *Shutter) {
		s.OnTerminating(f)
	}
}

var ErrShutterWasAlreadyDown = errors.New("saferun was called on an already-shutdown shutter")

// LockedInit allows you to run a function only if the shutter is not down yet,
// with the assurance that the it will not run its callback functions
// during the execution of your function.
//
// This is useful to prevent race conditions, where the func given to "LockedInit"
// should increase a counter and the func given to OnShutdown should decrease it.
//
// WARNING: never call Shutdown from within your LockedInit function,
// it will deadlock. Also, keep these init functions as short as
// possible.
//
// NOTE: This was previously named SafeRun
func (s *Shutter) LockedInit(fn func() error) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.IsTerminating() {
		return ErrShutterWasAlreadyDown
	}
	return fn()
}

func (s *Shutter) Shutdown(err error) {
	var execute = false
	s.once.Do(func() {
		execute = true
	})

	if !execute {
		return
	}

	s.lock.Lock()
	s.err = err
	close(s.terminatingCh)
	s.lock.Unlock()

	s.terminatingFuncLock.Lock()
	for _, call := range s.terminatingFunc {
		call(err)
	}
	s.terminatingFuncLock.Unlock()

	s.lock.Lock()
	// the err has been handle above thus it will be available
	close(s.terminatedCh)
	s.lock.Unlock()

	s.terminatedFuncLock.Lock()
	for _, call := range s.terminatedFunc {
		call(err)
	}
	s.terminatedFuncLock.Unlock()
}

func (s *Shutter) Terminating() <-chan struct{} {
	return s.terminatingCh
}

func (s *Shutter) IsTerminating() bool {
	select {
	case <-s.terminatingCh:
		return true
	default:
		return false
	}
}

func (s *Shutter) Terminated() <-chan struct{} {
	return s.terminatedCh
}

func (s *Shutter) IsTerminated() bool {
	select {
	case <-s.terminatedCh:
		return true
	default:
		return false
	}
}

// OnTerminating registers an additional handler to be triggered on
// `Shutdown()`. These calls will be blocking and will
// occur when the shutter is in the process of shutting down.
func (s *Shutter) OnTerminating(f func(error)) {
	s.terminatingFuncLock.Lock()
	s.terminatingFunc = append(s.terminatingFunc, f)
	s.terminatingFuncLock.Unlock()
}

// OnTerminated registers an additional handler to be triggered on
// `Shutdown()`. These calls will be blocking and will
// occur when the shutter has shutdown
func (s *Shutter) OnTerminated(f func(error)) {
	s.terminatedFuncLock.Lock()
	s.terminatedFunc = append(s.terminatedFunc, f)
	s.terminatedFuncLock.Unlock()
}

func (s *Shutter) Err() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.err
}
