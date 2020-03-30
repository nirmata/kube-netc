recv:
	go build -o recv -tags="linux_bpf" recv.go

promserv:
	go build -o promserv -tags="linux_bpf" promserv.go

clean:
	rm recv promserv
