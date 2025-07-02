build:
	go build -o bin/hit ./cmd/hit
	go build -o bin/server ./cmd/server

clean:
	rm bin/*

test:
	go test ./...