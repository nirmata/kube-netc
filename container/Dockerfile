FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN apk add build-base bcc linux-headers
RUN GOARCH=amd64 CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags="linux_bpf" -o main .

FROM busybox:latest
COPY --from=builder /build/ /app/
WORKDIR /app
ENTRYPOINT ["/bin/sh", "./container/entrypoint.sh"]
