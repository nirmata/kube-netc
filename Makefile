BUILD_ARGS="-tags=linux_bpf"
GIVE_SUDO=sudo -E env PATH=$(PATH)

recv:
	go build -o recv $(BUILD_ARGS) examples/recv.go

promserv:
	go build -o promserv $(BUILD_ARGS) examples/promserv.go

bps:
	go build -o bps $(BUILD_ARGS) examples/bps.go

tests:
	$(GIVE_SUDO) go test $(BUILD_ARGS) ./pkg/tracker 

clean:
	rm recv promserv bps
