package process

import (
	"fmt"
	"syscall"
	"time"

	"github.com/aphistic/sweet"
	"github.com/efritz/glock"
	. "github.com/onsi/gomega"

	"github.com/go-nacelle/config"
)

type WatcherSuite struct{}

func (s *WatcherSuite) TestNoErrors(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		watcher = newWatcher(errChan, outChan)
	)

	watcher.watch()

	// Nil errors do not go on out chan
	errChan <- errMeta{nil, makeNamedInitializer("a"), false}
	errChan <- errMeta{nil, makeNamedInitializer("b"), false}
	errChan <- errMeta{nil, makeNamedInitializer("c"), false}
	Consistently(outChan).ShouldNot(Receive())

	// Closing err chan should shutdown watcher
	close(errChan)
	Eventually(outChan).Should(BeClosed())

	// Ensure we unblock
	Eventually(outChan).Should(BeClosed())
}

func (s *WatcherSuite) TestFatalErrorBeginsShutdown(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		watcher = newWatcher(errChan, outChan)
	)

	watcher.watch()

	errChan <- errMeta{nil, makeNamedInitializer("a"), true}
	errChan <- errMeta{nil, makeNamedInitializer("b"), true}
	Consistently(outChan).ShouldNot(Receive())
	Consistently(watcher.shutdownSignal).ShouldNot(BeClosed())

	errChan <- errMeta{fmt.Errorf("utoh"), makeNamedInitializer("c"), true}
	Eventually(outChan).Should(Receive(MatchError("utoh")))
	Eventually(watcher.shutdownSignal).Should(BeClosed())
	Consistently(outChan).ShouldNot(Receive())

	// Additional errors
	errChan <- errMeta{nil, makeNamedInitializer("a"), true}
	errChan <- errMeta{nil, makeNamedInitializer("b"), true}
	Consistently(outChan).ShouldNot(Receive())

	// And the same behavior above applies
	close(errChan)
	Eventually(outChan).Should(BeClosed())
}

func (s *WatcherSuite) TestNilErrorBeginsShutdown(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		watcher = newWatcher(errChan, outChan)
	)

	watcher.watch()

	errChan <- errMeta{nil, makeNamedInitializer("a"), false}
	Consistently(outChan).ShouldNot(Receive())
	Eventually(watcher.shutdownSignal).Should(BeClosed())

	// Cleanup
	close(errChan)
	Eventually(outChan).Should(BeClosed())
}

func (s *WatcherSuite) TestSignals(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		watcher = newWatcher(errChan, outChan)
	)

	defer close(errChan)
	watcher.watch()

	// Try to ensure watcher is waiting on signal
	<-time.After(time.Millisecond * 100)

	// First signal
	syscall.Kill(syscall.Getpid(), shutdownSignals[0])
	Eventually(watcher.shutdownSignal).Should(BeClosed())

	// Second signal
	Consistently(watcher.abortSignal).ShouldNot(BeClosed())
	syscall.Kill(syscall.Getpid(), shutdownSignals[0])
	Eventually(watcher.abortSignal).Should(BeClosed())
}

func (s *WatcherSuite) TestExternalHaltRequestBeginsShutdown(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		watcher = newWatcher(errChan, outChan)
	)

	watcher.watch()

	watcher.halt()
	Eventually(watcher.shutdownSignal).Should(BeClosed())

	// Cleanup
	close(errChan)
	Eventually(outChan).Should(BeClosed())
}

func (s *WatcherSuite) TestShutdownTimeout(t sweet.T) {
	var (
		errChan = make(chan errMeta)
		outChan = make(chan error)
		clock   = glock.NewMockClock()
		watcher = newWatcher(
			errChan,
			outChan,
			withWatcherClock(clock),
			withWatcherShutdownTimeout(time.Second*10),
		)
	)

	defer close(errChan)
	watcher.watch()

	watcher.halt()
	Consistently(outChan).ShouldNot(BeClosed())
	clock.BlockingAdvance(time.Second * 10)
	Eventually(outChan).Should(BeClosed())
}

//
//

func makeNamedInitializer(name string) namedInitializer {
	initializer := InitializerFunc(func(config.Config) error {
		return nil
	})

	meta := newInitializerMeta(initializer)
	meta.name = name
	return meta
}
