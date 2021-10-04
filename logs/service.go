package logs

import (
	"context"
	"fmt"
	"github.com/aaronland/go-aws-session"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func GetServiceWithDSN(ctx context.Context, dsn string) (*cloudwatchlogs.CloudWatchLogs, error) {

	sess, err := session.NewSessionWithDSN(dsn)

	if err != nil {
		return nil, fmt.Errorf("Failed to create session, %w", err)
	}

	return GetServiceWithSession(ctx, sess)
}

func GetServiceWithSession(ctx context.Context, sess *aws_session.Session) (*cloudwatchlogs.CloudWatchLogs, error) {
	svc := cloudwatchlogs.New(sess)
	return svc, nil
}
