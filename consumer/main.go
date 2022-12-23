package main

import (
	"context"
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

	// client
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

	// attempt to consume data every 1 sec
	for {
		fmt.Println("-- keep scanning...")
		records, err := kinesisClient.GetRecords(context.TODO(), &kinesis.GetRecordsInput{
			ShardIterator: shardIterator,
		})
		if err != nil {
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		// process the data
		if len(records.Records) > 0 {
			log.Printf("GetRecords Data: %v\n", records.Records[0].Data)
		}
		shardIterator = records.NextShardIterator
		time.Sleep(1000 * time.Millisecond)
	}
}
