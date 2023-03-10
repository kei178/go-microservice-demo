# go-microservice-demo

![go-microservice-demo diagram](go-microservice-demo-diagram.jpeg)

## Consumer

Cosume data from the stream:

```
go run consumer/main.go --stream mystream
```


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

Restart the kinesis container:

```
docker-compose restart kinesis
```


### Useful Commands

Source: https://docs.aws.amazon.com/cli/latest/reference/kinesis/index.html

Create a stream:

```
AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis create-stream --stream-name mystream --shard-count 1
```


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

## Postgres

### Setup

Run Postgres in your local:

```
docker-compose up -d
```

Access psql console:

```
docker exec -it go-microservice-demo_postgres_1 bash
psql
```

Create a `testdb` database:

```
CREATE DATABASE testdb;
```

Create the `users` table:

```
CREATE TABLE users (
	id serial PRIMARY KEY,
	username VARCHAR ( 50 ) UNIQUE NOT NULL,
	email VARCHAR ( 255 ) UNIQUE NOT NULL,
	fullname VARCHAR ( 255 ) NOT NULL,
	created_at TIMESTAMP NOT NULL
);
```

### pgweb

You can access Postgres via pgweb at  `http://localhost:8088/`.
