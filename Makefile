recv:
	go build -o recv -tags="linux_bpf" recv.go

clean:
	rm recv
