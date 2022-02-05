all: build test package

build:
	GOOS=linux go build main.go

test: 
	GOOS=linux go test ./...

package:
	zip -o main.zip main

clean: 
	rm main