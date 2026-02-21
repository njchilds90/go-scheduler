// Package scheduler provides a lightweight, deterministic event scheduler
// inspired by Python's sched module.
//
// It allows scheduling functions to run at a specific time or after a delay,
// with priority-based ordering.
//
// This package is designed to be simple, embeddable, and useful for scripting,
// automation, testing, and AI agent orchestration.
package scheduler

import (
	"container/heap"
	"sync"
	"time"
)

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

// Scheduler executes scheduled events.
type Scheduler struct {
	mu     sync.Mutex
	events eventQueue
}

// New creates a new Scheduler.
func New() *Scheduler {
	s := &Scheduler{
		events: make(eventQueue, 0),
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
	return s.Schedule(time.Now().Add(delay), priority, action)
}

// Cancel removes a scheduled event.
func (s *Scheduler) Cancel(event *Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.index >= 0 {
		heap.Remove(&s.events, event.index)
	}
}

// Run executes events in order until no events remain.
func (s *Scheduler) Run() {
	for {
		s.mu.Lock()
		if s.events.Len() == 0 {
			s.mu.Unlock()
			return
		}

		next := s.events[0]
		now := time.Now()

		if now.Before(next.Time) {
			wait := next.Time.Sub(now)
			s.mu.Unlock()
			time.Sleep(wait)
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
