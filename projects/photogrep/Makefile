
all:

SRCS=photogrep.go csheet.go pan.go clickimg.go cache.go

photogrep: $(SRCS)
	go build $^

fmt:
	gofmt -w $(SRCS)

clean:
	rm -f photogrep
