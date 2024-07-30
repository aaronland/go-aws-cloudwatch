package logs

import (
	"context"
	"iter"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type FilterLogGroupFunc func(context.Context, *types.LogGroup) (bool, error)

func GetLogGroups(ctx context.Context, cl *cloudwatchlogs.Client, filters ...FilterLogGroupFunc) iter.Seq2[*types.LogGroup, error] {

	return func(yield func(*types.LogGroup, error) bool) {

		cursor := ""

		for {

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			opts := &cloudwatchlogs.DescribeLogGroupsInput{
				// LogGroupName: aws.String(log_group),
			}

			if cursor != "" {
				opts.NextToken = aws.String(cursor)
			}

			rsp, err := cl.DescribeLogGroups(ctx, opts)

			if err != nil {
				yield(nil, err)
				return
			}

			for _, s := range rsp.LogGroups {

				include_group := true

				for _, f := range filters {

					ok, err := f(ctx, &s)

					if err != nil {
						yield(nil, err)
						return
					}

					if !ok {
						include_group = false
						break
					}
				}

				if !include_group {
					continue
				}

				if !yield(&s, nil) {
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
