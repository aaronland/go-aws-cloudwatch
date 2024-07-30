// log-groups will emit the names of all the CloudWatch log groups for a given AWS account to STDOUT.
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

	flag.Parse()

	ctx := context.Background()

	cloudwatch_cl, err := logs.NewClient(ctx, *cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	for g, err := range logs.GetLogGroups(ctx, cloudwatch_cl) {

		if err != nil {
			log.Fatalf("Failed to get log groups, %v", err)
		}

		fmt.Println(*g.LogGroupName)
	}

}
