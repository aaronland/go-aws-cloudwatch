// log-stream-events will emit all of the events for a particular CloudWatch log stream to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-aws-cloudwatch/logs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"log"
)

func main() {

	cloudwatch_dsn := flag.String("cloudwatch-dsn", "", "A valid aaronland/go-aws-session DSN string.")
	cloudwatch_group := flag.String("log-group", "", "A valid CloudWatch log group name.")
	cloudwatch_stream := flag.String("log-stream", "", "A valid CloudWatch log stream name.")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_svc, err := logs.GetServiceWithDSN(ctx, *cloudwatch_dsn)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	event_ch := make(chan *cloudwatchlogs.OutputLogEvent)

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
