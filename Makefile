SOURCES=$(shell git ls-files . gen parse)
.PHONY: clean test test-unit test-integration

test: test-unit test-integration

test-integration: test/gomuti test/testify test/gomuti/generate test/testify/generate

test-unit: vendor
	ginkgo -r -skipPackage test

test/gomuti: install
	make clean
	mongoose -mock gomuti test/fixtures
	cd test/gomuti && ginkgo

test/gomuti/generate: install
	cd test/gomuti/generate && go generate
	cd test/gomuti/generate && ginkgo
	@echo "Checking for correct usage of -name flag:"
	cd test/gomuti/generate; grep -qv MockNotGenerated *.go

test/testify: install
	make clean
	mongoose -mock testify test/fixtures
	cd test/testify && ginkgo

test/testify/generate: install
	cd test/testify/generate && go generate
	cd test/testify/generate && ginkgo
	@echo "Checking for correct usage of -name flag:"
	cd test/testify/generate; grep -qv MockNotGenerated *.go

install: vendor $(SOURCES)
	go build -o mongoose *.go
	mv mongoose $$GOPATH/bin/mongoose

vendor: glide.yaml
	glide install

clean:
	rm -f test/fixtures/mock*.go
	rm -f test/gomuti/generate/mock*.go
	rm -f test/testify/generate/mock*.go
