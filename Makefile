SOURCES=$(shell git ls-files gen mock parse)
.PHONY: clean test test-unit test-integration

test: test-unit test-integration

test-integration: test/mongoose test/testify

test-unit:
	ginkgo -r -skipPackage test

test/mongoose: clean $(SOURCES)
	rm -f test/fixtures/mock*.go
	go run main.go test/fixtures
	cd test/mongoose && ginkgo

test/testify: clean $(SOURCES)
	rm -f test/fixtures/mock*.go
	go run main.go --mock=testify test/fixtures
	cd test/testify && ginkgo

clean:
	rm -f test/fixtures/mock*.go
