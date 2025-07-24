FROM ghcr.io/hpinc/krypton/krypton-go-builder as builder

ADD . /go/src/webhookworker
WORKDIR /go/src/webhookworker

# build the webhookworker binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w" \
    ./webhookworker.go

# Create the final minio image with the webhookworker configured.
FROM minio/minio:RELEASE.2020-06-22T03-12-50Z
RUN apk add curl jq openssl
COPY get_mc.sh /
RUN /get_mc.sh
COPY config_server.sh /
COPY start.sh /
COPY --from=builder /go/src/webhookworker/webhookworker /usr/local/bin
ENV MINIO_ACCESS_KEY=minioadmin
ENV MINIO_SECRET_KEY=minioadmin
ENV FS_BUCKET_NAMES="krypton_fs1,krypton_fs2"
ENV PUBLIC_BUCKET_NAMES=""
ENV SQS_HOST=sqs.local.test
ENV SQS_PORT=9324
ENV SQS_QUEUE_NAME=fs-notification
ENV WEBHOOK_PORT=9001
ENV SQS_RETRY_AFTER=3
EXPOSE 9000
ENTRYPOINT ["/start.sh"]
