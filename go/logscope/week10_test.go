package logscope

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWeek10Merge(t *testing.T) {
	defer catchTODO(t)
	destination := Summary{Valid: 1, Invalid: 2, Levels: map[Level]uint64{LevelInfo: 1}, Services: map[string]ServiceSummary{"api": {1, 5, 5}}}
	source := Summary{Valid: 2, Invalid: 1, Levels: map[Level]uint64{LevelError: 2}, Services: map[string]ServiceSummary{"api": {1, 2, 2}, "z": {1, 7, 7}}}
	if err := mergeSummaries(&destination, source); err != nil {
		t.Fatal(err)
	}
	if destination.Valid != 3 || destination.Invalid != 3 || destination.Levels[LevelInfo] != 1 || destination.Levels[LevelError] != 2 || destination.Services["api"] != (ServiceSummary{2, 7, 5}) || destination.Services["z"] != (ServiceSummary{1, 7, 7}) {
		t.Errorf("merge = %#v", destination)
	}
	if source.Valid != 2 || source.Services["api"] != (ServiceSummary{1, 2, 2}) {
		t.Error("merge mutated source")
	}
	source.Services["z"] = ServiceSummary{Count: 99}
	source.Levels[LevelError] = 99
	if destination.Services["z"].Count != 1 || destination.Levels[LevelError] != 2 {
		t.Error("merge retained source maps")
	}
	var zero Summary
	if err := mergeSummaries(&zero, Summary{}); err != nil {
		t.Errorf("merging zero summaries: %v", err)
	}
}

func TestWeek10MergeOverflow(t *testing.T) {
	for _, field := range []string{"valid", "invalid", "level", "service_count", "service_total"} {
		t.Run(field, func(t *testing.T) {
			defer catchTODO(t)
			destination := Summary{Levels: map[Level]uint64{}, Services: map[string]ServiceSummary{}}
			source := Summary{Valid: 1, Invalid: 1, Levels: map[Level]uint64{LevelInfo: 1}, Services: map[string]ServiceSummary{"api": {1, 1, 1}}}
			switch field {
			case "valid":
				destination.Valid = math.MaxUint64
			case "invalid":
				destination.Invalid = math.MaxUint64
			case "level":
				destination.Levels[LevelInfo] = math.MaxUint64
			case "service_count":
				destination.Services["api"] = ServiceSummary{Count: math.MaxUint64}
			case "service_total":
				destination.Services["api"] = ServiceSummary{TotalMS: math.MaxUint64}
			}
			if err := mergeSummaries(&destination, source); !errors.Is(err, ErrOverflow) {
				t.Errorf("merge overflow = %v; want ErrOverflow", err)
			}
		})
	}
}

func TestWeek10ParallelParity(t *testing.T) {
	inputs := map[string][]byte{
		"sample":     fixture(t, "sample.log"),
		"empty":      fixture(t, "empty.log"),
		"malformed":  fixture(t, "malformed.log"),
		"crlf":       []byte("INFO api 7\r\nERROR worker 9\r\n\r\n"),
		"fewer_jobs": []byte("WARN a 0"),
	}
	for name, input := range inputs {
		for workers := 1; workers <= 4; workers++ {
			t.Run(fmt.Sprintf("%s/%d", name, workers), func(t *testing.T) {
				defer catchTODO(t)
				sequential, err := Analyze(input)
				if err != nil {
					t.Fatal(err)
				}
				got, err := AnalyzeParallel(context.Background(), input, workers)
				if err != nil {
					t.Fatal(err)
				}
				if Render(got) != Render(sequential) {
					t.Errorf("parallel=%q; sequential=%q", Render(got), Render(sequential))
				}
			})
		}
	}
}

func TestWeek10Validation(t *testing.T) {
	defer catchTODO(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for _, workers := range []int{-1, 0, 5} {
		summary, err := AnalyzeParallel(ctx, nil, workers)
		if !errors.Is(err, ErrWorkers) {
			t.Errorf("workers=%d: %v; validation must precede cancellation", workers, err)
		}
		requireZeroSummary(t, summary)
	}
	summary, err := AnalyzeParallel(ctx, []byte("INFO a 1"), 1)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("pre-canceled context: %v", err)
	}
	requireZeroSummary(t, summary)
	var called atomic.Bool
	_, _ = analyzeParallel(ctx, nil, 1, func(context.Context, *Summary, []byte) error {
		called.Store(true)
		return nil
	})
	if called.Load() {
		t.Error("pre-canceled call started processing")
	}
	expired, stop := context.WithDeadline(context.Background(), time.Time{})
	defer stop()
	summary, err = AnalyzeParallel(expired, nil, 1)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expired context error = %v; want DeadlineExceeded", err)
	}
	requireZeroSummary(t, summary)
}

type poolOutcome struct {
	summary Summary
	err     error
}

func TestWeek10WorkerLifecycle(t *testing.T) {
	for _, workerFailure := range []bool{false, true} {
		t.Run(map[bool]string{false: "caller_cancellation", true: "worker_failure"}[workerFailure], func(t *testing.T) {
			const workers = 4
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			entered := make(chan struct{}, 64)
			releaseFailure := make(chan struct{})
			done := make(chan poolOutcome, 1)
			var active, maximum, exited atomic.Int32
			var mutex sync.Mutex
			locals := make(map[*Summary]bool)
			injected := errors.New("injected worker failure")
			process := func(workerContext context.Context, local *Summary, line []byte) error {
				current := active.Add(1)
				defer active.Add(-1)
				defer exited.Add(1)
				for previous := maximum.Load(); current > previous; previous = maximum.Load() {
					if maximum.CompareAndSwap(previous, current) {
						break
					}
				}
				mutex.Lock()
				first := len(locals) == 0
				locals[local] = true
				mutex.Unlock()
				entered <- struct{}{}
				if workerFailure && first {
					select {
					case <-releaseFailure:
						return injected
					case <-workerContext.Done():
						return workerContext.Err()
					}
				}
				<-workerContext.Done()
				return workerContext.Err()
			}
			go func() {
				outcome := poolOutcome{}
				defer func() {
					if value := recover(); value != nil {
						outcome.err = fmt.Errorf("pool panicked: %v", value)
					}
					done <- outcome
				}()
				outcome.summary, outcome.err = analyzeParallel(ctx, []byte(strings.Repeat("INFO a 1\n", 40)), workers, process)
			}()
			guard := time.NewTimer(10 * time.Second)
			defer guard.Stop()
			for index := 0; index < workers; index++ {
				select {
				case <-entered:
				case outcome := <-done:
					t.Fatalf("returned before all workers entered: %v", outcome.err)
				case <-guard.C:
					cancel()
					t.Fatal("deadlock guard: workers never entered")
				}
			}
			if workerFailure {
				close(releaseFailure)
			} else {
				cancel()
			}
			select {
			case outcome := <-done:
				cause := error(context.Canceled)
				if workerFailure {
					cause = injected
				}
				if !errors.Is(outcome.err, cause) {
					t.Errorf("pool error = %v; want %v", outcome.err, cause)
				}
				requireZeroSummary(t, outcome.summary)
				if maximum.Load() != workers || active.Load() != 0 || exited.Load() != workers {
					t.Errorf("max=%d active=%d exited=%d; want 4,0,4 (all callbacks joined)", maximum.Load(), active.Load(), exited.Load())
				}
				mutex.Lock()
				localCount := len(locals)
				mutex.Unlock()
				if localCount != workers {
					t.Errorf("workers shared local summaries: %d distinct, want %d", localCount, workers)
				}
				if workerFailure && ctx.Err() != nil {
					t.Error("pool canceled the caller's context")
				}
			case <-guard.C:
				cancel()
				t.Fatal("deadlock guard: cancellation did not join all workers")
			}
		})
	}
}

func TestWeek10WorkerOverflow(t *testing.T) {
	defer catchTODO(t)
	// One worker receives both records: no scheduling-dependent distribution.
	summary, err := analyzeParallel(context.Background(), []byte("one\ntwo\n"), 1,
		func(ctx context.Context, local *Summary, line []byte) error {
			if string(line) == "one" {
				local.Valid = math.MaxUint64
			} else {
				return ErrOverflow
			}
			return nil
		})
	if !errors.Is(err, ErrOverflow) {
		t.Errorf("worker aggregate overflow = %v", err)
	}
	requireZeroSummary(t, summary)
}

func TestWeek10PoolMergeOverflow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	done := make(chan poolOutcome, 1)
	go func() {
		outcome := poolOutcome{}
		defer func() {
			if value := recover(); value != nil {
				outcome.err = fmt.Errorf("pool panicked: %v", value)
			}
			done <- outcome
		}()
		outcome.summary, outcome.err = analyzeParallel(ctx, []byte("one\ntwo\n"), 2,
			func(ctx context.Context, local *Summary, line []byte) error {
				local.Valid = math.MaxUint64
				entered <- struct{}{}
				select {
				case <-release:
					return nil
				case <-ctx.Done():
					return ctx.Err()
				}
			})
	}()
	guard := time.NewTimer(10 * time.Second)
	defer guard.Stop()
	for index := 0; index < 2; index++ {
		select {
		case <-entered:
		case outcome := <-done:
			t.Fatalf("returned before both independent summaries: %v", outcome.err)
		case <-guard.C:
			t.Fatal("deadlock guard: merge-overflow callbacks never entered")
		}
	}
	close(release)
	select {
	case outcome := <-done:
		if !errors.Is(outcome.err, ErrOverflow) {
			t.Errorf("pool merge error = %v; want ErrOverflow", outcome.err)
		}
		requireZeroSummary(t, outcome.summary)
	case <-guard.C:
		t.Fatal("deadlock guard: pool failed to join before checked merge")
	}
}

func TestWeek10CancellationBeforeSuccess(t *testing.T) {
	defer catchTODO(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	summary, err := analyzeParallel(ctx, []byte("INFO a 1\n"), 1,
		func(ctx context.Context, local *Summary, line []byte) error {
			cancel()
			return nil
		})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("caller canceled at final record: %v; want context.Canceled", err)
	}
	requireZeroSummary(t, summary)
}

func TestWeek10NoSourceMutation(t *testing.T) {
	defer catchTODO(t)
	source := Summary{Levels: map[Level]uint64{LevelInfo: 2}, Services: map[string]ServiceSummary{"a": {2, 9, 8}}}
	var destination Summary
	if err := mergeSummaries(&destination, source); err != nil {
		t.Fatal(err)
	}
	destination.Levels[LevelInfo]++
	destination.Services["a"] = ServiceSummary{}
	if !reflect.DeepEqual(source.Levels, map[Level]uint64{LevelInfo: 2}) || source.Services["a"] != (ServiceSummary{2, 9, 8}) {
		t.Error("destination aliases source storage")
	}
}
