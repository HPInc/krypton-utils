#!/bin/sh

STORAGE_NAME=krypton-fs
HOST=localhost
PORT=9000
SERVER="http://$HOST:$PORT"
HTTP_FORBIDDEN=403
WEBHOOK_HOST=localhost
WEBHOOK_PORT=9001

wait_for_server() {
	this_host=${1:-$HOST}
	this_port=${2:-$PORT}
	# wait for server start
	while ! nc -z "$this_host" "$this_port"; do
		echo "Waiting for service start..."
		sleep 1
	done
	echo "service started"

	# check server access
	result=0
	while [ "$result" -ne "$HTTP_FORBIDDEN" ]; do
		echo "Waiting for service ready status. Result = $result"
		result=$(curl -so /dev/null --write-out "%{http_code}" "$SERVER")
		sleep 1
	done
}

setup_buckets() {
	# add alias
	mc config host add \
	"$STORAGE_NAME" "$SERVER" \
	"$MINIO_ACCESS_KEY" \
	"$MINIO_SECRET_KEY"

	# add buckets
	for i in $(echo "$FS_BUCKET_NAMES" | sed 's/,/ /g'); do
		mc mb --ignore-existing "$STORAGE_NAME/$i"
	done
	for i in $(echo "$PUBLIC_BUCKET_NAMES" | sed 's/,/ /g'); do
		mc policy set public "$STORAGE_NAME/$i"
	done
}

restart_service() {
	mc admin service restart "$STORAGE_NAME"
	wait_for_server "$HOST" "$PORT"
}

# set up notification to sqs. this is setup via a local webhook
# which forwards to sqs
setup_notification() {
	webhookworker &
	# wait for webhook server
	wait_for_server "$WEBHOOK_HOST" "$WEBHOOK_PORT"

	# create queue dir for cache. prevents webhook overload, allows 429 retries
	mkdir /tmp/queue

	# set up storage notification
	mc admin config \
	set "$STORAGE_NAME" \
	notify_webhook:1 endpoint="http://$WEBHOOK_HOST:$WEBHOOK_PORT" queue_dir="/tmp/queue"

	# notification setup requires a restart of storage service
	restart_service
	echo "restart done"
}

# wait for minio storage service
wait_for_server "$HOST" "$PORT"

# set up buckets
setup_buckets
echo "setup buckets done"

# set up notification
setup_notification
echo "setup notification done"

# setup notification for buckets
for i in $(echo "$FS_BUCKET_NAMES" | sed 's/,/ /g'); do
	mc event add "$STORAGE_NAME/$i" arn:minio:sqs::1:webhook --event put
done
