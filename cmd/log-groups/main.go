// log-groups will emit the names of all the CloudWatch log groups for a given AWS account to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/aaronland/go-aws-cloudwatch/logs"
)

func main() {

	var cloudwatch_uri string
	var verbose bool

	flag.StringVar(&cloudwatch_uri, "cloudwatch-uri", "", "A valid aaronland/go-aws-auth URI.")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	cloudwatch_cl, err := logs.NewClient(ctx, cloudwatch_uri)

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
