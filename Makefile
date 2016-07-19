SOURCES=$(shell git ls-files . gen parse)
.PHONY: clean test test-unit test-integration

test: test-unit test-integration

test-integration: test/gomuti test/testify

test-unit:
	ginkgo -r -skipPackage test

test/gomuti: $(SOURCES)
	make clean
	go run main.go -mock gomuti test/fixtures
	cd test/gomuti && ginkgo

test/testify: $(SOURCES)
	make clean
	go run main.go -mock testify test/fixtures
	cd test/testify && ginkgo

clean:
	rm -f test/fixtures/mock*.go
