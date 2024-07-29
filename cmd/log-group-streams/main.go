// log-group-streams will emit the names of all the log streams in a given CloudWatch log group to STDOUT.
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

	cloudwatch_loggroup := flag.String("log-group", "", "A valid CloudWatch log group name.")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_svc, err := logs.NewClient(ctx, *cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	streams, err := logs.GetLogGroupStreams(ctx, cloudwatch_svc, *cloudwatch_loggroup)

	if err != nil {
		log.Fatalf("Failed to get log streams, %v", err)
	}

	for _, s := range streams {
		fmt.Println(*s.LogStreamName)
	}

}
