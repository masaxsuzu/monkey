test:
	go vet ./...
	goimports -l ./
	go test ./...
run:
	go run main.go

tojs:
	gopherjs build playground/main.go -o docs\playground.js -o docs/playground.js