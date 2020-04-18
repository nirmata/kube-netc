FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN ls
RUN apk add build-base bcc linux-headers
RUN GOARCH=amd64 CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags="linux_bpf" -o main .

FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]