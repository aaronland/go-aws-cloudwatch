package logs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type FilterLogGroupFunc func(context.Context, *types.LogGroup) (bool, error)

func GetLogGroups(ctx context.Context, cl *cloudwatchlogs.Client, filters ...FilterLogGroupFunc) ([]*types.LogGroup, error) {

	groups := make([]*types.LogGroup, 0)

	cursor := ""

	for {

		select {
		case <-ctx.Done():
			return groups, nil
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
			return nil, fmt.Errorf("Failed to describe groups, %w", err)
		}

		for _, s := range rsp.LogGroups {

			include_group := true

			for _, f := range filters {

				ok, err := f(ctx, &s)

				if err != nil {
					return nil, fmt.Errorf("Filter func for %s failed, %w", *s.LogGroupName, err)
				}

				if !ok {
					ok = false
					break
				}
			}

			if !include_group {
				continue
			}

			groups = append(groups, &s)
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

	return groups, nil
}
