test:
	go fmt ./...
	go test ./...
validate:
	go fmt ./...
	go vet ./...
	goimports -l ./
	go test ./...
bench:
	go test ./benchmark -bench Run -cpu 1
run:
	go run main.go

tojs:
	gopherjs build playground/main.go -o docs\playground.js -o docs/playground.js