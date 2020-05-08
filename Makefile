test:
	@go test ./...

fmt:
	@go fmt ./...

check:
	@echo "Checking formatting"
	@[ -z "$(shell gofmt -l ./..)" ]
	@echo "Vetting"
	@go vet
	@echo "Running tests"
	@make test
