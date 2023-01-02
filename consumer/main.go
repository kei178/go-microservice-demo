package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"

	_ "github.com/lib/pq"
)

type UserPayload struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func main() {
	var (
		streamName      = flag.String("stream", "", "Stream name")
		kinesisEndpoint = flag.String("endpoint", "http://localhost:4567", "Kinesis endpoint")
		awsRegion       = flag.String("region", "us-west-2", "AWS Region")
	)
	flag.Parse()

	db := connectToDB()
	defer db.Close()

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
				var payload UserPayload
				err := json.Unmarshal([]byte(r.Data), &payload)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Printf("GetRecords Data: %v\n", payload)
				fmt.Println("---")

				err = insertUser(db, payload)
				if err != nil {
					fmt.Printf("failed to save user payload: %v", err)
				}
			}
		}
		shardIterator = resp.NextShardIterator
		time.Sleep(interval)
	}
}

const (
	dbhost     = "localhost"
	dbport     = 5555
	dbuser     = "root"
	dbpassword = "password"
	dbname     = "testdb"
)

func connectToDB() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpassword, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Panic(err)
	}

	// check db
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Connected to Postgres!")
	return db
}

func insertUser(db *sql.DB, payload UserPayload) error {
	_, err := db.Exec(
		`INSERT INTO users ("username", "email", "fullname", "created_at") VALUES ($1, $2, $3, $4)`,
		payload.Username,
		payload.Email,
		payload.Firstname+" "+payload.Lastname,
		time.Now(),
	)
	return err
}
