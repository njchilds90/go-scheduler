# go-scheduler

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-scheduler.svg)](https://pkg.go.dev/github.com/njchilds90/go-scheduler)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/go-scheduler)](https://goreportcard.com/report/github.com/njchilds90/go-scheduler)

A lightweight, deterministic, priority-based event scheduler for Go inspired by Python’s `sched` module.

Designed for scripting, automation, deterministic testing, and AI agent orchestration.

---

## Why This Library Exists

Go provides:

- `time.Timer`
- `time.Ticker`
- `time.After`

But Go does **not** provide:

- A structured priority-based event scheduler
- A deterministic single-threaded event loop
- A simple embeddable task scheduler
- A minimal scripting-friendly scheduling API

Python has `sched`.  
Go does not have a clean equivalent.

`go-scheduler` fills that gap.

---

## Features

- Schedule tasks at a specific time
- Schedule tasks after a delay
- Priority-based ordering
- Deterministic execution
- Cancel scheduled events
- Zero external dependencies
- Minimal, standard-library style API
- AI-agent friendly design

---

## Installation

```bash
go get github.com/njchilds90/go-scheduler
Quick Start
package main

import (
	"fmt"
	"time"

	"github.com/njchilds90/go-scheduler"
)

func main() {
	s := scheduler.New()

	s.ScheduleAfter(2*time.Second, 1, func() {
		fmt.Println("Executed after 2 seconds")
	})

	s.RunBlocking()
}
Scheduling at a Specific Time
s := scheduler.New()

runAt := time.Now().Add(5 * time.Second)

s.Schedule(runAt, 1, func() {
	fmt.Println("Executed at specific time")
})

s.Run()
Priority Ordering

If two events are scheduled for the same time, lower priority values execute first.

now := time.Now()

s.Schedule(now, 2, func() {
	fmt.Println("Priority 2")
})

s.Schedule(now, 1, func() {
	fmt.Println("Priority 1")
})

s.Run()

Output:

Priority 1
Priority 2
Cancel an Event
event := s.ScheduleAfter(5*time.Second, 1, func() {
	fmt.Println("This will not run")
})

s.Cancel(event)
Deterministic Execution Model

go-scheduler runs in a single goroutine and executes events sequentially.

This makes it ideal for:

Deterministic testing

Embedded scripting engines

Game loops

AI agent decision cycles

Time-based simulations

Retry/backoff systems

No worker pools.
No background goroutines.
No hidden concurrency.

AI Agent Use Cases

This scheduler is especially useful for AI systems:

Action planning with delayed execution

Coordinating multi-step agent workflows

Time-based decision making

Controlled simulation environments

Retry logic with backoff

Event-driven orchestration loops

AI systems benefit from deterministic event ordering — this library provides that.

Design Philosophy

Small and focused

No dependencies

Predictable behavior

Clean API

Inspired by Python, written the Go way

Easy to understand and audit

API Overview
Create Scheduler
s := scheduler.New()
Schedule at Time
func (s *Scheduler) Schedule(at time.Time, priority int, action func()) *Event
Schedule After Delay
func (s *Scheduler) ScheduleAfter(delay time.Duration, priority int, action func()) *Event
Cancel Event
func (s *Scheduler) Cancel(event *Event)
Run Scheduler
func (s *Scheduler) Run()
func (s *Scheduler) RunBlocking()
Versioning

This project follows Semantic Versioning.

Example:

go get github.com/njchilds90/go-scheduler@v0.1.0
Roadmap

Potential future additions:

Context cancellation support

Non-blocking async runner

Simulated/test clock

Event introspection API

Metrics hooks

GitHub Actions CI

100% coverage badge

Contributing

Fork the repository

Create a feature branch

Submit a Pull Request

Keep it simple. Keep it deterministic.

License

MIT License

Copyright (c) 2026 Nicholas Childs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction...

Author

Nicholas Childs
GitHub: https://github.com/njchilds90
