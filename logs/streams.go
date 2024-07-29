package logs

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type FilterLogStreamFunc func(context.Context, *types.LogStream) (bool, error)

func LogsStreamsWithBytes(ctx context.Context, s *types.LogStream) (bool, error) {

	if *s.StoredBytes == 0 {
		return false, nil
	}

	return true, nil
}

func LogStreamsSinceFunc(ctx context.Context, ts int64) FilterLogStreamFunc {

	fn := func(ctx context.Context, s *types.LogStream) (bool, error) {

		return *s.LastEventTimestamp >= ts, nil
	}

	return fn
}

func GetMostRecentStreamForLogGroup(ctx context.Context, cl *cloudwatchlogs.Client, log_group string) (*types.LogStream, error) {

	filters := []FilterLogStreamFunc{
		LogsStreamsWithBytes,
	}

	streams, err := GetLogGroupStreams(ctx, cl, log_group, filters...)

	if err != nil {
		return nil, fmt.Errorf("Failed to determine streams for log group, %w", err)
	}

	for _, s := range streams {

		log.Println(*s.LogStreamName, *s.FirstEventTimestamp, *s.LastEventTimestamp)
	}

	count := len(streams)

	return streams[count-1], nil
}

func GetLogGroupStreams(ctx context.Context, cl *cloudwatchlogs.Client, log_group string, filters ...FilterLogStreamFunc) ([]*types.LogStream, error) {

	streams := make([]*types.LogStream, 0)

	cursor := ""

	for {

		select {
		case <-ctx.Done():
			return streams, nil
		default:
			// pass
		}

		opts := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(log_group),
		}

		if cursor != "" {
			opts.NextToken = aws.String(cursor)
		}

		rsp, err := cl.DescribeLogStreams(ctx, opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to describe streams for %s, %w", log_group, err)
		}

		for _, s := range rsp.LogStreams {

			include_stream := true

			for _, f := range filters {

				ok, err := f(ctx, &s)

				if err != nil {
					return nil, fmt.Errorf("Filter func for %s failed, %w", *s.LogStreamName, err)
				}

				if !ok {
					ok = false
					break
				}
			}

			if !include_stream {
				continue
			}

			streams = append(streams, &s)
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

	return streams, nil
}
