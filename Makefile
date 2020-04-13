.DEFAULT_GOAL := build

IMAGENAME := kube-netc
LD_FLAGS:=-ldflags="-w -s"
BUILD_ARGS := -tags="linux_bpf"
GIVE_SUDO := sudo -E env PATH=$(PATH):$(GOPATH)

recv:
	go build -o recv $(BUILD_ARGS) examples/recv.go

promserv:
	go build -o promserv $(BUILD_ARGS) examples/promserv.go

bps:
	go build -o bps $(BUILD_ARGS) examples/bps.go

tests:
	$(GIVE_SUDO) go test $(BUILD_ARGS) ./pkg/tracker 

build:
	go get ./pkg/...
	go build $(BUILD_ARGS) -o main main.go

build-docker:
	GOARCH=amd64 CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo $(LD_FLAGS) $(BUILD_ARGS) -o main .
	docker build -t $(IMAGENAME) -f Dockerfile .

run-docker:
	$(GIVE_SUDO) docker run --rm -v /sys/kernel/debug:/sys/kernel/debug -v /sys/fs/cgroup:/sys/fs/cgroup -v /sys/fs/bpf:/sys/fs/bpf --privileged $(IMAGENAME)

run: build-docker run-docker

clean:
	rm -f recv promserv bps main
