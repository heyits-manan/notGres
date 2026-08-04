.PHONY: build test vet run clean

build:
	go build -o bin/notgres .

test:
	go test ./...

vet:
	go vet ./...

run:
	go run . --data-dir ./data --port 5432

clean:
	rm -rf bin