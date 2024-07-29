// empty-streams will list all the CloudWatch log streams with 0 stored bytes and optionally
// remove them.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aaronland/go-aws-cloudwatch/logs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func main() {

	cloudwatch_uri := flag.String("cloudwatch-uri", "", "...")

	prune := flag.Bool("prune", false, "Remove log streams with no events.")
	dryrun := flag.Bool("dryrun", false, "Go through the motions but don't actually remove any log streams.")

	max_workers := flag.Int("max-workers", 100, "The maximum number of concurrent workers.")

	flag.Parse()

	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, *max_workers)

	for i := 0; i < *max_workers; i++ {
		throttle <- true
	}

	limiter := time.Tick(200 * time.Millisecond)

	cloudwatch_cl, err := logs.NewClient(ctx, *cloudwatch_uri)

	if err != nil {
		log.Fatalf("Failed to create service, %v", err)
	}

	groups, err := logs.GetLogGroups(ctx, cloudwatch_cl)

	if err != nil {
		log.Fatalf("Failed to get log groups, %v", err)
	}

	stream_filter := func(ctx context.Context, s *types.LogStream) (bool, error) {

		if *s.StoredBytes == 0 {
			return true, nil
		}

		return false, nil
	}

	wg := new(sync.WaitGroup)

	for _, g := range groups {

		fmt.Println(*g.LogGroupName)

		streams, err := logs.GetLogGroupStreams(ctx, cloudwatch_cl, *g.LogGroupName, stream_filter)

		if err != nil {
			log.Fatalf("Failed to get log streams for %s, %v", *g.LogGroupName, err)
		}

		for _, s := range streams {

			name := fmt.Sprintf("%s#%s\n", *g.LogGroupName, *s.LogStreamName)
			fmt.Println(name)

			if *prune {

				events_opts := &logs.GetLogEventsOptions{
					LogGroupName:  *g.LogGroupName,
					LogStreamName: *s.LogStreamName,
				}

				events, err := logs.GetLogEventsAsList(ctx, cloudwatch_cl, events_opts)

				if err != nil {
					log.Fatalf("Failed to get events for %s (%s), %v", *g.LogGroupName, *s.LogStreamName, err)
				}

				if len(events) > 0 {
					log.Printf("%s (%s) has events (%d) even though stored bytes is 0\n", *g.LogGroupName, *s.LogStreamName, len(events))
					continue
				}

				wg.Add(1)

				go func(g *types.LogGroup, s *types.LogStream) {

					defer func() {
						throttle <- true
						wg.Done()
					}()

					<-throttle
					<-limiter

					n := fmt.Sprintf("%s#%s\n", *g.LogGroupName, *s.LogStreamName)

					if *dryrun {
						log.Printf("Prune %s (dryrun)\n", n)
					}

					err := pruneStream(ctx, cloudwatch_cl, g, s)

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

func pruneStream(ctx context.Context, cl *cloudwatchlogs.Client, g *types.LogGroup, s *types.LogStream) error {

	opts := &cloudwatchlogs.DeleteLogStreamInput{
		LogGroupName:  g.LogGroupName,
		LogStreamName: s.LogStreamName,
	}

	_, err := cl.DeleteLogStream(ctx, opts)
	return err
}
