package cloudwatch

import (
	"context"
	"github.com/aaronland/go-aws-session"
	"github.com/aws/aws-sdk-go/aws"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	_ "log"
)

type GetLogEventsOptions struct {
	LogGroupName    string
	LogStreamName   string
	StartFromHead   bool
	LogEventChannel chan *cloudwatchlogs.OutputLogEvent
}

func GetLogEventsWithDSN(ctx context.Context, dsn string, opts *GetLogEventsOptions) error {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return err
	}

	return GetLogEventsWithSession(ctx, sess, opts)
}

func GetLogEventsWithSession(ctx context.Context, sess *aws_session.Session, opts *GetLogEventsOptions) error {

	svc := cloudwatchlogs.New(sess)
	return GetLogEventsWithService(ctx, svc, opts)
}

func GetLogEventsAsListWithService(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, opts *GetLogEventsOptions) ([]*cloudwatchlogs.OutputLogEvent, error) {

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

	err := GetLogEventsWithService(ctx, svc, opts)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetLogEventsWithService(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, opts *GetLogEventsOptions) error {

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
