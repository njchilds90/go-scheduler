package scheduler

import (
	"sync"
	"testing"
	"time"
)

func TestScheduleAfter(t *testing.T) {
	s := New()

	var mu sync.Mutex
	executed := false

	s.ScheduleAfter(100*time.Millisecond, 1, func() {
		mu.Lock()
		executed = true
		mu.Unlock()
	})

	s.Run()

	mu.Lock()
	defer mu.Unlock()

	if !executed {
		t.Fatal("expected scheduled function to execute")
	}
}

func TestPriorityOrder(t *testing.T) {
	s := New()

	var result []int
	var mu sync.Mutex

	now := time.Now()

	s.Schedule(now, 2, func() {
		mu.Lock()
		result = append(result, 2)
		mu.Unlock()
	})

	s.Schedule(now, 1, func() {
		mu.Lock()
		result = append(result, 1)
		mu.Unlock()
	})

	s.Run()

	if len(result) != 2 {
		t.Fatal("expected 2 results")
	}

	if result[0] != 1 {
		t.Fatal("priority ordering failed")
	}
}
