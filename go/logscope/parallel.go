package logscope

import "context"

// mergeSummaries does not change source or retain its maps. Discard destination
// after an error. Check every count and sum; maxima are selected, not added.
func mergeSummaries(destination *Summary, source Summary) error {
	// TODO G10.1: Merge uint64 aggregates with checkedAdd and independent maps.
	return ErrUnimplemented
}

type lineProcessor func(context.Context, *Summary, []byte) error

// analyzeParallel is the learner pool and the deterministic lifecycle test seam.
// process is called synchronously within a worker, never in another goroutine.
func analyzeParallel(ctx context.Context, input []byte, workers int, process lineProcessor) (Summary, error) {
	// TODO G10.2: Validate workers before ctx.Err, derive cancellation, start
	// exactly workers workers, produce jobs here, close jobs, collect every
	// result, then merge. See the spec for error precedence and channel bounds.
	return Summary{}, ErrUnimplemented
}

// AnalyzeParallel's supplied adapter uses the same record operation as Analyze.
func AnalyzeParallel(ctx context.Context, input []byte, workers int) (Summary, error) {
	return analyzeParallel(ctx, input, workers, func(ctx context.Context, summary *Summary, line []byte) error {
		return accumulateLine(summary, line)
	})
}
