package logs

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	_ "log"
)

type FilterLogGroupFunc func(context.Context, *cloudwatchlogs.LogGroup) (bool, error)

func GetLogGroups(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, filters ...FilterLogGroupFunc) ([]*cloudwatchlogs.LogGroup, error) {

	groups := make([]*cloudwatchlogs.LogGroup, 0)

	cursor := ""

	for {

		opts := &cloudwatchlogs.DescribeLogGroupsInput{
			// LogGroupName: aws.String(log_group),
		}

		if cursor != "" {
			opts.NextToken = aws.String(cursor)
		}

		rsp, err := svc.DescribeLogGroups(opts)

		if err != nil {
			return nil, fmt.Errorf("Failed to describe groups, %w", err)
		}

		for _, s := range rsp.LogGroups {

			include_group := true

			for _, f := range filters {

				ok, err := f(ctx, s)

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

			groups = append(groups, s)
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
