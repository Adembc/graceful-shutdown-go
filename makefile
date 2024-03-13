run:
	go run main.go
build:
	go build -o bin/server main.go
build-run:
	go build -o bin/server main.go
	./bin/server