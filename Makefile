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

full-test:
	@echo "<html><head/></html>" > index.html
	@go build
	@REACT_APP_TEST_VAR=XYZ ./runtime-js-env
	@cat index.html
