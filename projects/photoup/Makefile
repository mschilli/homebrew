all:

SRCS=photoup.go util.go tmpl.go image.go

photoup: $(SRCS)
	go build $^

fmt:
	gofmt -w $(SRCS)

clean:
	rm -f photoup
