// Package scheduler provides a lightweight, deterministic, priority-based
// event scheduler inspired by Python's sched module.
//
// It allows scheduling functions to run at a specific time or after a delay,
// with priority-based ordering.
//
// This package is designed to be simple, embeddable, and useful for scripting,
// automation, testing, and AI agent orchestration.
//
// New Features:
//   - Context cancellation support
//   - Non-blocking RunAsync()
//   - Simulated clock support for deterministic testing
//   - Graceful Stop()
//   - Pending event inspection
//
// The scheduler runs in a single goroutine and executes events sequentially.
package scheduler

import (
	"container/heap"
	"context"
	"sync"
	"time"
)

//
// CLOCK ABSTRACTION
//

// Clock defines the time source used by the Scheduler.
// This enables real-time or simulated time operation.
type Clock interface {
	Now() time.Time
	Sleep(d time.Duration)
}

// RealClock uses the system clock.
type RealClock struct{}

// Now returns the current time.
func (RealClock) Now() time.Time { return time.Now() }

// Sleep pauses execution for the specified duration.
func (RealClock) Sleep(d time.Duration) { time.Sleep(d) }

// SimulatedClock is a controllable clock for deterministic testing.
type SimulatedClock struct {
	mu  sync.Mutex
	now time.Time
}

// NewSimulatedClock creates a simulated clock starting at the given time.
func NewSimulatedClock(start time.Time) *SimulatedClock {
	return &SimulatedClock{now: start}
}

// Now returns the simulated current time.
func (c *SimulatedClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

// Sleep advances the simulated time.
func (c *SimulatedClock) Sleep(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

// Advance manually advances simulated time.
func (c *SimulatedClock) Advance(d time.Duration) {
	c.Sleep(d)
}

//
// EVENT & PRIORITY QUEUE
//

// Event represents a scheduled function call.
type Event struct {
	Time     time.Time
	Priority int
	Action   func()
	index    int
}

// eventQueue implements heap.Interface for Events.
type eventQueue []*Event

func (eq eventQueue) Len() int { return len(eq) }

func (eq eventQueue) Less(i, j int) bool {
	if eq[i].Time.Equal(eq[j].Time) {
		return eq[i].Priority < eq[j].Priority
	}
	return eq[i].Time.Before(eq[j].Time)
}

func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].index = i
	eq[j].index = j
}

func (eq *eventQueue) Push(x interface{}) {
	n := len(*eq)
	event := x.(*Event)
	event.index = n
	*eq = append(*eq, event)
}

func (eq *eventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	event := old[n-1]
	event.index = -1
	*eq = old[0 : n-1]
	return event
}

//
// SCHEDULER
//

// Scheduler executes scheduled events sequentially.
type Scheduler struct {
	mu      sync.Mutex
	events  eventQueue
	clock   Clock
	running bool
	stopCh  chan struct{}
}

// New creates a new Scheduler using the real system clock.
func New() *Scheduler {
	return NewWithClock(RealClock{})
}

// NewWithClock creates a Scheduler with a custom Clock.
// Useful for testing with SimulatedClock.
func NewWithClock(clock Clock) *Scheduler {
	s := &Scheduler{
		events: make(eventQueue, 0),
		clock:  clock,
		stopCh: make(chan struct{}),
	}
	heap.Init(&s.events)
	return s
}

// Schedule schedules an action to run at a specific time.
func (s *Scheduler) Schedule(at time.Time, priority int, action func()) *Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	event := &Event{
		Time:     at,
		Priority: priority,
		Action:   action,
	}

	heap.Push(&s.events, event)
	return event
}

// ScheduleAfter schedules an action to run after a delay.
func (s *Scheduler) ScheduleAfter(delay time.Duration, priority int, action func()) *Event {
	return s.Schedule(s.clock.Now().Add(delay), priority, action)
}

// Cancel removes a scheduled event.
func (s *Scheduler) Cancel(event *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.index >= 0 {
		heap.Remove(&s.events, event.index)
	}
}

// Pending returns the number of scheduled events.
func (s *Scheduler) Pending() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.events.Len()
}

// Stop gracefully stops the scheduler loop.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if s.running {
		close(s.stopCh)
		s.running = false
	}
	s.mu.Unlock()
}

// Run executes events in order until no events remain.
func (s *Scheduler) Run() {
	s.runWithContext(context.Background())
}

// RunWithContext executes scheduled events until completion
// or until the context is canceled.
func (s *Scheduler) RunWithContext(ctx context.Context) {
	s.runWithContext(ctx)
}

// internal run logic
func (s *Scheduler) runWithContext(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopCh = make(chan struct{})
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.running = false
		s.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		default:
		}

		s.mu.Lock()
		if s.events.Len() == 0 {
			s.mu.Unlock()
			return
		}

		next := s.events[0]
		now := s.clock.Now()

		if now.Before(next.Time) {
			wait := next.Time.Sub(now)
			s.mu.Unlock()
			s.clock.Sleep(wait)
			continue
		}

		heap.Pop(&s.events)
		s.mu.Unlock()

		next.Action()
	}
}

// RunBlocking runs the scheduler in the current goroutine.
func (s *Scheduler) RunBlocking() {
	s.Run()
}

// RunAsync runs the scheduler in a new goroutine.
func (s *Scheduler) RunAsync() {
	go s.Run()
}
