// log-groups will emit the names of all the CloudWatch log groups for a given AWS account to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-aws-cloudwatch/logs"
	"log"
)

func main() {

	cloudwatch_dsn := flag.String("cloudwatch-dsn", "", "A valid aaronland/go-aws-session DSN string.")

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
