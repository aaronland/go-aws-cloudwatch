package logs

import (
	"context"
	"iter"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type FilterLogStreamFunc func(context.Context, *types.LogStream) (bool, error)

func LogsStreamsWithBytesFunc() FilterLogStreamFunc {

	fn := func(ctx context.Context, s *types.LogStream) (bool, error) {

		if *s.StoredBytes == 0 {
			return false, nil
		}

		return true, nil
	}

	return fn
}

func LogStreamsSinceFunc(ts int64) FilterLogStreamFunc {

	fn := func(ctx context.Context, s *types.LogStream) (bool, error) {
		return *s.LastEventTimestamp >= ts, nil
	}

	return fn
}

func GetMostRecentStreamForLogGroup(ctx context.Context, cl *cloudwatchlogs.Client, log_group string) (*types.LogStream, error) {

	filters := []FilterLogStreamFunc{
		LogsStreamsWithBytesFunc(),
	}

	var recent *types.LogStream

	for s, err := range GetLogGroupStreams(ctx, cl, log_group, filters...) {

		if err != nil {
			return nil, err
		}

		recent = s
		break
	}

	return recent, nil
}

func GetLogGroupStreamsSince(ctx context.Context, cl *cloudwatchlogs.Client, log_group string, ts int64) iter.Seq2[*types.LogStream, error] {

	filters := []FilterLogStreamFunc{
		LogsStreamsWithBytesFunc(),
		LogStreamsSinceFunc(ts),
	}

	return GetLogGroupStreams(ctx, cl, log_group, filters...)
}

func GetLogGroupStreams(ctx context.Context, cl *cloudwatchlogs.Client, log_group string, filters ...FilterLogStreamFunc) iter.Seq2[*types.LogStream, error] {

	return func(yield func(*types.LogStream, error) bool) {

		cursor := ""

		for {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs#DescribeLogStreamsInput
			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs#DescribeLogStreamsOutput
			// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs@v1.37.3/types#LogStream

			opts := &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: aws.String(log_group),
				// Default to most recent
				Descending: aws.Bool(true),
				OrderBy:    types.OrderByLastEventTime,
			}

			if cursor != "" {
				opts.NextToken = aws.String(cursor)
			}

			rsp, err := cl.DescribeLogStreams(ctx, opts)

			if err != nil {
				slog.Error("Failed to describe log stream", "error", err)
				yield(nil, err)
				return
			}

			for _, s := range rsp.LogStreams {

				include_stream := true

				for _, f := range filters {

					ok, err := f(ctx, &s)

					if err != nil {
						slog.Error("Stream filter failed, skipping", "error", err)
						if !yield(nil, err) {
							return
						}
					}

					if !ok {
						include_stream = false
						break
					}
				}

				if !include_stream {
					continue
				}

				if !yield(&s, err) {
					return
				}
			}

			if rsp.NextToken == nil {
				break
			}

			if *rsp.NextToken != "" && *rsp.NextToken != cursor {
				cursor = *rsp.NextToken
			} else {
				break
			}

		}
	}
}
