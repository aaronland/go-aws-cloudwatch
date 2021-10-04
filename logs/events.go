package logs

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type GetLogEventsOptions struct {
	LogGroupName    string
	LogStreamName   string
	StartFromHead   bool
	LogEventChannel chan *cloudwatchlogs.OutputLogEvent
}

func GetLogEventsAsList(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, opts *GetLogEventsOptions) ([]*cloudwatchlogs.OutputLogEvent, error) {

	events := make([]*cloudwatchlogs.OutputLogEvent, 0)
	events_ch := make(chan *cloudwatchlogs.OutputLogEvent)
	done_ch := make(chan bool)

	defer func() {
		done_ch <- true
	}()

	go func() {

		for {
			select {
			case <-done_ch:
				return
			case e := <-events_ch:
				events = append(events, e)
			default:
				// pass
			}
		}

	}()

	opts.LogEventChannel = events_ch

	err := GetLogEvents(ctx, svc, opts)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetLogEvents(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, opts *GetLogEventsOptions) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var cursor string

	for {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		req := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(opts.LogGroupName),
			LogStreamName: aws.String(opts.LogStreamName),
			StartFromHead: aws.Bool(opts.StartFromHead),
		}

		if cursor != "" {
			req.NextToken = aws.String(cursor)
		}

		rsp, err := svc.GetLogEvents(req)

		if err != nil {
			return err
		}

		if len(rsp.Events) == 0 {
			break
		}

		for _, e := range rsp.Events {
			opts.LogEventChannel <- e
		}

		// sigh... (20190213/thisisaaronland)

		if *rsp.NextForwardToken != "" && *rsp.NextForwardToken != cursor {
			cursor = *rsp.NextForwardToken
		} else {
			break
		}

	}

	return nil
}
