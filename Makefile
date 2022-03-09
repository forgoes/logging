

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out -coverpkg=./... ./...
	# go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=coverage.out -o coverage.html