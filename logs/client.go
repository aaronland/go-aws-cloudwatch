package logs

import (
	"context"
	_ "fmt"

	"github.com/aaronland/go-aws-auth"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func NewClient(ctx context.Context, uri string) (*cloudwatchlogs.Client, error) {

	cfg, err := auth.NewConfig(ctx, uri)

	if err != nil {
		return nil, err
	}

	return cloudwatchlogs.NewFromConfig(cfg), nil

}
