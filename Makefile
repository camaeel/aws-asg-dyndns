all: build test package

build:
	GOOS=linux go build -o main

mock:
	GOOS=linux go generate ./...

test: mock
	GOOS=linux go test ./...

package:
	zip -o main.zip main

clean: clean_main clean_mocks

clean_main:
	rm main

clean_mocks:
	find . -name '*_mocks.go' -delete
