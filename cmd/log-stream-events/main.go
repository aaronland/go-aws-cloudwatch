// log-stream-events will emit all of the events for a particular CloudWatch log stream to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aaronland/go-aws-cloudwatch/logs"
)

func main() {

	cloudwatch_uri := flag.String("cloudwatch-uri", "", "...")

	cloudwatch_group := flag.String("log-group", "", "A valid CloudWatch log group name.")
	cloudwatch_stream := flag.String("log-stream", "", "A valid CloudWatch log stream name.")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_cl, err := logs.NewClient(ctx, *cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	opts := &logs.GetLogEventsOptions{
		LogGroupName:  *cloudwatch_group,
		LogStreamName: *cloudwatch_stream,
	}

	for e, err := range logs.GetLogEvents(ctx, cloudwatch_cl, opts) {

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(e)
	}

}
