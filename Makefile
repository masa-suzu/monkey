test:
	go vet ./...
	golint ./...
	go test ./...
run:
	go run main.go