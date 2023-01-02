#!/bin/bash

JQ_CHECK=$(which jq)
if [ -z "$JQ_CHECK" ]; then
  echo
  echo "This script requires the jq JSON processor. Please install for your OS from https://stedolan.github.io/jq/download/"
  echo
  exit 1
fi

if [ $# -ne 1 ]; then
  echo
  echo "usage: $0 <stream_name>"
  echo
  exit 1
fi

# Set the stream name
STREAM_NAME=$1

# Choose the iterator type:
# TRIM HORIZON is for starting at the begining of the Kinesis Stream.
# This can take a while if you have a lot of records.
# To use TRIM HORIZON, uncomment the following line:
# TYPE=TRIM_HORIZON

# AT_TIMESTAMP allows you to go back to a point in time. This is set for going back one hour
# To use AT_TIMESTAMP, uncomment the following two lines:
# TIMESTAMP=$(($(date +%s) - 3600)).000
# TYPE="AT_TIMESTAMP --timestamp $TIMESTAMP"

# LATEST means start at the most current point in the stream and read forward
TYPE=TRIM_HORIZON

echo "-- Fetching SHARD_IDS"

# Get a list of shards
SHARD_IDS=$(AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis describe-stream --stream-name $STREAM_NAME | jq -r .StreamDescription.Shards[].ShardId)

echo $SHARD_IDS

echo "-- Fetching SHARD_ITERATORS"

# Get all the starting points
for shard_id in $SHARD_IDS ; do
  echo "SHARD_ITERATOR for shard_id ${shard_id}"
  echo $(AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-shard-iterator --shard-id $shard_id --shard-iterator-type $TYPE --stream-name $STREAM_NAME --query 'ShardIterator')
done
