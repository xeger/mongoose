SOURCES=$(shell git ls-files gen mock parse)
.PHONY: clean test test-unit test-integration

test: $(SOURCES)

test-integration: test/mongoose test/testify

test-unit:
	ginkgo -r -skipPackage test

test/mongoose: $(SOURCES)
	go run main.go test/mongoose
	cd test/mongoose && ginkgo

test/testify: $(SOURCES)
	go run main.go --mock=testify test/testify
	cd test/testify && ginkgo

clean:
	rm -f test/mongoose/mock*.go
	rm -f test/testify/mock*.go
