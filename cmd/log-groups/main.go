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

	// prefix := flag.String("prefix", "", "...")

	flag.Parse()

	ctx := context.Background()

	cloudwatch_svc, err := logs.GetServiceWithDSN(ctx, *cloudwatch_dsn)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	groups, err := logs.GetLogGroups(ctx, cloudwatch_svc)

	if err != nil {
		log.Fatalf("Failed to get log groups, %v", err)
	}

	for _, g := range groups {
		fmt.Println(*g.LogGroupName)
	}

}
