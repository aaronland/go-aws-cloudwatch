package logs

import (
	"context"
	"iter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type GetLogEventsOptions struct {
	LogGroupName  string
	LogStreamName string
	StartFromHead bool
}

func GetLogEvents(ctx context.Context, cl *cloudwatchlogs.Client, opts *GetLogEventsOptions) iter.Seq2[*types.OutputLogEvent, error] {

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

			req := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(opts.LogGroupName),
				LogStreamName: aws.String(opts.LogStreamName),
				StartFromHead: aws.Bool(opts.StartFromHead),
			}

			if cursor != "" {
				req.NextToken = aws.String(cursor)
			}

			rsp, err := cl.GetLogEvents(ctx, req)

			if err != nil {

				if !yield(nil, err) {
					return
				}
			}

			if len(rsp.Events) == 0 {
				break
			}

			for _, e := range rsp.Events {

				if !yield(&e, nil) {
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
