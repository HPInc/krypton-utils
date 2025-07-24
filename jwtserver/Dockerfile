FROM ghcr.io/hpinc/krypton/krypton-go-builder AS builder

ADD . /go/src/krypton-utils/jwtserver
WORKDIR /go/src/krypton-utils/jwtserver

# build the source
RUN make build-binaries

# use a minimal alpine image
FROM ghcr.io/hpinc/krypton/krypton-go-base

# set working directory
WORKDIR /go/bin

COPY --from=builder /go/src/krypton-utils/jwtserver/bin/jwtserver .

EXPOSE 9090/tcp

# run the binary
CMD ["./jwtserver"]
