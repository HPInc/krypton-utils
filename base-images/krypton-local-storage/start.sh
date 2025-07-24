#!/bin/sh
/config_server.sh &
minio server /data
