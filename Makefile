helm-docs:
	cd cmd/helm-docs && go build
	mv cmd/helm-docs/helm-docs .

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm helm-docs
