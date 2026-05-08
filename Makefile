.PHONY: test test-race test-vet test-cover

test:
	go test ./...

test-race:
	go test -race ./...

test-vet:
	go vet ./...

test-cover:
	go test ./... -coverprofile coverage.out
	go tool cover -func coverage.out
