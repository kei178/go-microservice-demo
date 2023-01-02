# go-microservice-demo


## Producer

Produce data to the stream:

```
cat producer/users.txt  | go run producer/main.go --stream mystream
```

## Kinesis

### Setup

Run [Kinesis Lite](https://github.com/mhart/kinesalite) in your local:

```
docker-compose up -d
```

### Useful Commands

Source: https://docs.aws.amazon.com/cli/latest/reference/kinesis/index.html

Get a list of streams:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis list-streams
```

Get info (`shard-id`) of a stream:


```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis describe-stream --stream-name mystream
```

Get a `shard-iterator` from a stream:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-shard-iterator --shard-id shardId-000000000000 --shard-iterator-type TRIM_HORIZON --stream-name mystream --query 'ShardIterator'
```

Get records with a `shard-iterator`:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-records --shard-iterator [shard-iterator]
```

Get the decoded first record from a `shard-iterator`:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-records --shard-iterator [shard-iterator] | jq -r '.Records[0].Data' | base64 --decode
```

### Scripts

Get a list of `shard-id` and `shard-iterator` from a stream:

```
bash scripts/get-shards.sh mystream
```

Get a list of records from a stream with a `shard-iterator`:

```
SHARD_ITERATOR=[shard-iterator]
bash scripts/get-records.sh mystream $SHARD_ITERATOR
```
