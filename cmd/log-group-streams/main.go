// log-group-streams will emit the names of all the log streams in a given CloudWatch log group to STDOUT.
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

	cloudwatch_loggroup := flag.String("log-group", "", "A valid CloudWatch log group name.")

	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	cloudwatch_svc, err := logs.NewClient(ctx, cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	for s, err := range logs.GetLogGroupStreams(ctx, cloudwatch_svc, *cloudwatch_loggroup) {

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(*s.LogStreamName, err)
	}

}
