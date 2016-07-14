SOURCES=$(shell git ls-files . gen parse)

test: test/mongoose test/testify

test/mongoose: $(SOURCES)
	go run main.go test/mongoose
	cd test/mongoose && ginkgo

test/testify: $(SOURCES)
	go run main.go --mock=testify test/testify
	cd test/testify && ginkgo

clean:
	rm -f test/mongoose/mock*.go
	rm -f test/testify/mock*.go
