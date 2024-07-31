package logs

import (
	"context"
	"iter"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// See also: https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html

type FilterLogEventFunc func(context.Context, *types.OutputLogEvent) (bool, error)

func FilterLambdaStartEndEventFunc() FilterLogEventFunc {

	fn := func(ctx context.Context, ev *types.OutputLogEvent) (bool, error) {

		if strings.HasPrefix(*ev.Message, "START RequestId: ") {
			return false, nil
		}

		if strings.HasPrefix(*ev.Message, "END RequestId: ") {
			return false, nil
		}

		return true, nil
	}

	return fn
}

type GetLogEventsOptions struct {
	LogGroupName  string
	LogStreamName string
	StartFromHead bool
	StartTime     int64
	EndTime       int64
	Filters       []FilterLogEventFunc
}

func GetLogEvents(ctx context.Context, cl *cloudwatchlogs.Client, opts *GetLogEventsOptions) iter.Seq2[*types.OutputLogEvent, error] {

	if opts.LogStreamName == "" {

		slog.Debug("Stream name is empty, polling all streams for group", "group", opts.LogGroupName)

		return func(yield func(*types.OutputLogEvent, error) bool) {

			for s, err := range GetLogGroupStreams(ctx, cl, opts.LogGroupName) {

				if err != nil {
					slog.Error("Failed to get streams for group", "group", opts.LogGroupName, "error", err)
					yield(nil, err)
					return
				}

				ev_opts := &GetLogEventsOptions{
					LogGroupName:  opts.LogGroupName,
					LogStreamName: *s.LogStreamName,
					StartFromHead: opts.StartFromHead,
					StartTime:     opts.StartTime,
					EndTime:       opts.EndTime,
					Filters:       opts.Filters,
				}

				slog.Debug("Get events for stream", "group", opts.LogGroupName, "stream", *s.LogStreamName, "first event", *s.FirstEventTimestamp, "last event", *s.LastEventTimestamp, "start from head", opts.StartFromHead)

				for ev, err := range GetLogEvents(ctx, cl, ev_opts) {

					if !yield(ev, err) {
						return
					}
				}

				if opts.StartTime > 0 {

					last := *s.LastEventTimestamp
					stop := opts.StartTime >= last

					// slog.Debug("Test timestamps", "start", opts.StartTime, "last", last, "stop", stop)

					if stop {
						break
					}
				}

			}
		}
	}

	return func(yield func(*types.OutputLogEvent, error) bool) {

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var cursor string

		for {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs#GetLogEventsInput

			req := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(opts.LogGroupName),
				LogStreamName: aws.String(opts.LogStreamName),
			}

			if opts.StartFromHead {
				req.StartFromHead = aws.Bool(opts.StartFromHead)
			}

			if opts.StartTime > 0 {
				req.StartTime = aws.Int64(opts.StartTime)
			}

			if opts.EndTime > 0 {
				req.EndTime = aws.Int64(opts.EndTime)
			}

			if cursor != "" {
				req.NextToken = aws.String(cursor)
			}

			slog.Debug("Get events", "group", opts.LogGroupName, "stream", opts.LogStreamName, "start", opts.StartTime, "cursor", cursor)

			rsp, err := cl.GetLogEvents(ctx, req)

			if err != nil {

				slog.Error("Failed to retrieve log events", "group", opts.LogGroupName, "stream", opts.LogStreamName, "error", err)
				if !yield(nil, err) {
					return
				}
			}

			if len(rsp.Events) == 0 {
				// slog.Debug("No events from group/stream", "group", opts.LogGroupName, "stream", opts.LogStreamName)
				break
			}

			for _, ev := range rsp.Events {

				filter_ok := true

				for _, f := range opts.Filters {

					ok, err := f(ctx, &ev)

					if err != nil {
						if !yield(nil, err) {
							return
						}
					}

					if !ok {
						filter_ok = false
						break

					}
				}

				if !filter_ok {
					continue
				}

				if !yield(&ev, nil) {
					return
				}
			}

			// sigh... (20190213/thisisaaronland)

			if *rsp.NextForwardToken != "" && *rsp.NextForwardToken != cursor {
				cursor = *rsp.NextForwardToken
			} else {
				break
			}

		}
	}

}
