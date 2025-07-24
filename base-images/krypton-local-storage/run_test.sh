#!/bin/sh
DOCKER_IMAGE_NAME=krypton-local-storage
FS_PORT=9000
FS_TEST_NAME=fs_storage_test
FS_TEST_USER=fstestadmin
FS_TEST_PASS=fstestpass
FS_TEST_BUCKET_PREFIX=fs-test
FS_TEST_BUCKETS="fs-test1,fs-test2,fs-test3"
FS_STORAGE_NAME=krypton-fs
FS_MAX_RETRIES=10
FS_DELAY_SECONDS=2

docker run \
	-p$FS_PORT:$FS_PORT \
	-eMINIO_ACCESS_KEY=$FS_TEST_USER \
	-eMINIO_SECRET_KEY=$FS_TEST_PASS \
	-eFS_BUCKET_NAMES=$FS_TEST_BUCKETS \
	-eSQS_HOST=172.19.0.2 \
	-eSQS_QUEUE_NAME=pending-enroll \
	--name $FS_TEST_NAME \
	--rm -d $DOCKER_IMAGE_NAME

i=0
bucket_count=$(echo $FS_TEST_BUCKETS | sed 's|,|\n|g' | grep -c $FS_TEST_BUCKET_PREFIX)
while [ "$i" -lt "$FS_MAX_RETRIES" ]; do
	count=$(docker exec $FS_TEST_NAME mc ls $FS_STORAGE_NAME | grep -c $FS_TEST_BUCKET_PREFIX)
	if [ "$count" -eq "$bucket_count" ]; then
		echo "Local storage configuration successfully verified"
		exit 0
	fi
	echo "Waiting for buckets to be configured. $i/$FS_MAX_RETRIES"
	sleep "$FS_DELAY_SECONDS"
	i=$((i+1))
done
exit 1
