# go-microservice-demo


## Producer

Produce data to the stream:

```
cat producer/users.txt  | go run producer/main.go --stream mystream
```

## Kinesis

### Setup

Run Kinesislite in your local:

```
docker-compose up -d
```

### Useful Commands

Sourece: https://docs.aws.amazon.com/cli/latest/reference/kinesis/index.html

Get a list of streams:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis list-streams
```

Get info of a stream:


```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis describe-stream --stream-name mystream
```

Get a shard-iterator from a stream:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-shard-iterator --shard-id [shard-id] --shard-iterator-type TRIM_HORIZON --stream-name mystream --query 'ShardIterator'
```

Get records with a shard iterator:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-records --shard-iterator [shard-iterator]
```
