// log-stream-events will emit all of the events for a particular CloudWatch log stream to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"

	"github.com/aaronland/go-aws-cloudwatch/logs"
	"github.com/sfomuseum/go-flags/multi"
)

func main() {

	var cloudwatch_uri string
	var cloudwatch_group string
	var cloudwatch_stream string
	var verbose bool

	var str_filters multi.MultiString

	flag.StringVar(&cloudwatch_uri, "cloudwatch-uri", "", "...")

	flag.StringVar(&cloudwatch_group, "log-group", "", "A valid CloudWatch log group name.")
	flag.StringVar(&cloudwatch_stream, "log-stream", "", "A valid CloudWatch log stream name.")

	// start / end stuff here

	flag.BoolVar(&verbose, "verbose", false, "")
	flag.Var(&str_filters, "filter", "")

	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	cloudwatch_cl, err := logs.NewClient(ctx, cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	opts := &logs.GetLogEventsOptions{
		LogGroupName:  cloudwatch_group,
		LogStreamName: cloudwatch_stream,
	}

	if len(str_filters) > 0 {

		filters := make([]logs.FilterLogEventFunc, 0)

		for _, str_f := range str_filters {

			switch str_f {
			case "lambda":
				filters = append(filters, logs.FilterLambdaStartEndEventFunc())
			default:
				slog.Error("Invalid or unsupported string filter", "filter", str_f)
			}
		}

		opts.Filters = filters
	}

	for ev, err := range logs.GetLogEvents(ctx, cloudwatch_cl, opts) {

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(*ev.Message)
	}

}
