// empty-streams will list all the CloudWatch log streams with 0 stored bytes and optionally
// remove them. IMPORTANT: It is not clear to me that this tool is working as expected yet...
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-aws-cloudwatch/logs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"log"
	"sync"
	"time"
)

func main() {

	cloudwatch_dsn := flag.String("cloudwatch-dsn", "region=us-west-2 credentials=session", "...")
	prune := flag.Bool("prune", false, "...")

	max_workers := flag.Int("max-workers", 100, "...")

	flag.Parse()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, *max_workers)

	for i := 0; i < *max_workers; i++ {
		throttle <- true
	}

	limiter := time.Tick(200 * time.Millisecond)

	cloudwatch_svc, err := logs.GetServiceWithDSN(ctx, *cloudwatch_dsn)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	groups, err := logs.GetLogGroups(ctx, cloudwatch_svc)

	if err != nil {
		log.Fatalf("Failed to get log groups, %v", err)
	}

	stream_filter := func(ctx context.Context, s *cloudwatchlogs.LogStream) (bool, error) {

		if *s.StoredBytes == 0 {
			return true, nil
		}

		return false, nil
	}

	wg := new(sync.WaitGroup)

	for _, g := range groups {

		fmt.Println(*g.LogGroupName)

		streams, err := logs.GetLogGroupStreams(ctx, cloudwatch_svc, *g.LogGroupName, stream_filter)

		if err != nil {
			log.Fatalf("Failed to get log streams for %s, %v", *g.LogGroupName, err)
		}

		for _, s := range streams {

			name := fmt.Sprintf("%s#%s\n", *g.LogGroupName, *s.LogStreamName)
			fmt.Println(name)

			if *prune {

				wg.Add(1)

				go func(g *cloudwatchlogs.LogGroup, s *cloudwatchlogs.LogStream) {

					defer func() {
						throttle <- true
						wg.Done()
					}()

					<-throttle
					<-limiter

					n := fmt.Sprintf("%s#%s\n", *g.LogGroupName, *s.LogStreamName)

					err := pruneStream(ctx, cloudwatch_svc, g, s)

					if err != nil {
						log.Println("Failed to remove %s (%s), %v", n, err)
					} else {
						log.Printf("Pruned %s\n", n)
					}
				}(g, s)
			}
		}

	}

	wg.Wait()
}

func pruneStream(ctx context.Context, svc *cloudwatchlogs.CloudWatchLogs, g *cloudwatchlogs.LogGroup, s *cloudwatchlogs.LogStream) error {

	opts := &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  g.LogGroupName,
		LogStreamName: s.LogStreamName,
	}

	_, err := svc.DeleteLogStream(opts)
	return err
}
