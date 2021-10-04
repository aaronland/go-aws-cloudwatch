package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-aws-cloudwatch/logs"
	"log"
)

func main() {

	cloudwatch_dsn := flag.String("cloudwatch-dsn", "region=us-west-2 credentials=session", "...")
	cloudwatch_loggroup := flag.String("cloudwatch-loggroup", "", "...")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_svc, err := logs.GetServiceWithDSN(ctx, *cloudwatch_dsn)

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
