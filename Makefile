BUILD_ARGS="-tags=linux_bpf"

recv:
	go build -o recv $(BUILD_ARGS) examples/recv.go

promserv:
	go build -o promserv $(BUILD_ARGS) examples/promserv.go

tests:
	go test $(BUILD_ARGS) ./pkg/tracker 

clean:
	rm examples/recv examples/promserv
