package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func main() {
	var (
		streamName      = flag.String("stream", "", "Stream name")
		kinesisEndpoint = flag.String("endpoint", "http://localhost:4567", "Kinesis endpoint")
		awsRegion       = flag.String("region", "us-west-2", "AWS Region")
	)
	flag.Parse()

	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           *kinesisEndpoint,
			SigningRegion: *awsRegion,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(*awsRegion),
		config.WithEndpointResolver(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("user", "pass", "token")),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	var kinesisClient = kinesis.NewFromConfig(cfg)

	// retrieve shard IDs
	streams, err := kinesisClient.DescribeStream(context.TODO(), &kinesis.DescribeStreamInput{StreamName: streamName})
	if err != nil {
		log.Fatalf("failed to fetch shard IDs: %v", err)
	}

	// retrieve shard iterators
	iteratorOutput, err := kinesisClient.GetShardIterator(context.TODO(), &kinesis.GetShardIteratorInput{
		ShardId:           aws.String(*streams.StreamDescription.Shards[0].ShardId),
		ShardIteratorType: "TRIM_HORIZON",
		StreamName:        streamName,
	})
	if err != nil {
		log.Fatalf("failed to fetch shard iterators: %v", err)
	}
	shardIterator := iteratorOutput.ShardIterator

	// attempt to consume data with a 2-sec interval
	var interval = 2000 * time.Millisecond
	for {
		fmt.Println("Keep scanning...")
		resp, err := kinesisClient.GetRecords(context.TODO(), &kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			time.Sleep(interval)
			continue
		}

		// process the data
		if len(resp.Records) > 0 {
			for _, r := range resp.Records {
				var result map[string]interface{}
				err := json.Unmarshal([]byte(r.Data), &result)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("GetRecords Data: %v\n", result)
				fmt.Println("---")
			}
		}
		shardIterator = resp.NextShardIterator
		time.Sleep(interval)
	}
}
