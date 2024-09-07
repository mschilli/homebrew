all:

gitmeta: gitmeta.go gmf.go parser.go gmf-entry.go gmf-ghuser.go gmf-ssh.go clone.go
	go build $^

clean:
	rm -f gitmeta

fmt:
	gofmt -w *.go
