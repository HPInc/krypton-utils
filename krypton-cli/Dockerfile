FROM ghcr.io/hpinc/krypton/krypton-go-builder AS builder

ADD . /go/src/cli
WORKDIR /go/src/cli

# build the source
RUN make vet go_build

# use a minimal alpine image
FROM ghcr.io/hpinc/krypton/krypton-go-base

# install common tools for self-contained runs
RUN apk add bash curl jq

# make krypton-cli available in path
COPY --from=builder /go/src/cli/bin/krypton-cli /usr/local/bin/krypton-cli
COPY --from=builder /go/src/cli/config/config.yaml /root/.krypton-cli/config.yaml

# Create a lower privileged user account to run cli
RUN adduser --disabled-password --gecos "" --home "/home/cliacct" \
    --shell "/sbin/nologin" --no-create-home --uid 1001 cliacct && \
    mkdir -p /tmp/krypton && \
    chown -R cliacct:cliacct /tmp/krypton

# run the binary
ENTRYPOINT ["krypton-cli"]
