package logs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"log"
)

type FilterLogStreamFunc func(context.Context, *cloudwatchlogs.LogStream) (bool, error)

func LogsStreamsWithBytes(ctx context.Context, s *cloudwatchlogs.LogStream) (bool, error) {

	if *s.StoredBytes == 0 {
		return false, nil
	}

	return true, nil
}

func GetMostRecentStreamForLogGroup(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, log_group string) (*cloudwatchlogs.LogStream, error) {

	filters := []FilterLogStreamFunc{
		LogsStreamsWithBytes,
	}

	streams, err := GetLogGroupStreams(ctx, svc, log_group, filters...)

	if err != nil {
		return nil, fmt.Errorf("Failed to determine streams for log group, %w", err)
	}

	for _, s := range streams {

		log.Println(*s.LogStreamName, *s.FirstEventTimestamp, *s.LastEventTimestamp)
	}

	count := len(streams)

	return streams[count-1], nil
}

func GetLogGroupStreams(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, log_group string, filters ...FilterLogStreamFunc) ([]*cloudwatchlogs.LogStream, error) {

	streams := make([]*cloudwatchlogs.LogStream, 0)

	cursor := ""

	for {

		opts := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(log_group),
		}

		if cursor != "" {
			opts.NextToken = aws.String(cursor)
		}

		rsp, err := svc.DescribeLogStreams(opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to describe streams for %s, %w", log_group, err)
		}

		for _, s := range rsp.LogStreams {

			include_stream := true

			for _, f := range filters {

				ok, err := f(ctx, s)

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

			streams = append(streams, s)
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
