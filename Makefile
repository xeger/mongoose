test: fixtures/testify
	go run main.go fixtures/testify
	cd fixtures/testify && ginkgo

clean:
	rm -f fixtures/testify/mock*.go
