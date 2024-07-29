// log-stream-events will emit all of the events for a particular CloudWatch log stream to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aaronland/go-aws-cloudwatch/logs"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func main() {

	cloudwatch_uri := flag.String("cloudwatch-uri", "", "...")

	cloudwatch_group := flag.String("log-group", "", "A valid CloudWatch log group name.")
	cloudwatch_stream := flag.String("log-stream", "", "A valid CloudWatch log stream name.")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_svc, err := logs.NewClient(ctx, *cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	event_ch := make(chan *types.OutputLogEvent)

	go func() {

		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-event_ch:
				fmt.Println(*ev.Message)
			}
		}
	}()

	opts := &logs.GetLogEventsOptions{
		LogGroupName:    *cloudwatch_group,
		LogStreamName:   *cloudwatch_stream,
		LogEventChannel: event_ch,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = logs.GetLogEvents(ctx, cloudwatch_svc, opts)

	if err != nil {
		log.Fatalf("Failed to get log events, %v")
	}

}
