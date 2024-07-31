// log-stream-events will emit all of the events for a particular CloudWatch log stream to STDOUT.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aaronland/go-aws-cloudwatch/logs"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/sfomuseum/iso8601duration"
)

func main() {

	var cloudwatch_uri string
	var cloudwatch_group string
	var cloudwatch_stream string
	var verbose bool

	var stdout bool
	var stderr bool
	var null bool

	var since string
	var str_filters multi.MultiString

	flag.StringVar(&cloudwatch_uri, "cloudwatch-uri", "", "A valid aaronland/go-aws-auth URI.")
	flag.StringVar(&cloudwatch_group, "log-group", "", "A valid CloudWatch log group name.")
	flag.StringVar(&cloudwatch_stream, "log-stream", "", "A valid CloudWatch log stream name.")

	flag.StringVar(&since, "since", "", "A valid ISO8061 duration string.")

	// start / end stuff here

	flag.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")
	flag.Var(&str_filters, "filter", "")

	flag.BoolVar(&stdout, "stdout", false, "Write log events to STDOUT.")
	flag.BoolVar(&stderr, "stderr", false, "Write log events to STDERR.")
	flag.BoolVar(&null, "null", false, "Write log events to /dev/null.")

	flag.Parse()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ctx := context.Background()

	writers := make([]io.Writer, 0)

	if stdout {
		writers = append(writers, os.Stdout)
	}

	if stderr {
		writers = append(writers, os.Stderr)
	}

	if null {
		writers = append(writers, io.Discard)
	}

	mw := io.MultiWriter(writers...)

	cloudwatch_cl, err := logs.NewClient(ctx, cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create client, %v", err)
	}

	opts := &logs.GetLogEventsOptions{
		LogGroupName:  cloudwatch_group,
		LogStreamName: cloudwatch_stream,
	}

	if since != "" {

		dur, err := duration.FromString(since)

		if err != nil {
			log.Fatalf("Failed to parse -since flag, %w", err)
		}

		now := time.Now()
		then := now.Add(-dur.ToDuration())

		opts.StartTime = then.Unix() * 1000

		slog.Debug("Filter events starting on or after", "dt", then.Format(time.RFC3339), "timestamp", then.Unix(), "start", opts.StartTime)
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

		m := strings.TrimSpace(*ev.Message)
		m = fmt.Sprintf("%s\n", m)
		
		mw.Write([]byte(m))		
	}

}
