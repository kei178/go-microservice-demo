JQ_CHECK=$(which jq)
if [ -z "$JQ_CHECK" ]; then
  echo
  echo "This script requires the jq JSON processor. Please install for your OS from https://stedolan.github.io/jq/download/"
  echo
  exit 1
fi

# Set the stream name
STREAM_NAME=$1
SHARD_ITERATOR=$2

# Start getting events from a shard and display them
DATA=$(AWS_ACCESS_KEY_ID=x AWS_SECRET_ACCESS_KEY=x aws --endpoint-url http://localhost:4567/ kinesis get-records --limit 50 --shard-iterator $SHARD_ITERATOR)
ROWS=$(echo $DATA | jq -r .Records[].Data?)
for row in $ROWS; do
  echo $row | base64 -d | jq .
done
