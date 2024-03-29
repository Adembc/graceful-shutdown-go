run:
	go run main.go
build:
	go build -o bin/server main.go
build-run:
	go build -o bin/server main.go
	./bin/server
calc:
	curl  'http://localhost:8080/calculation'