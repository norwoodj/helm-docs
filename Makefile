helm-docs:
	go build github.com/norwoodj/helm-docs/cmd/helm-docs

install:
	go install github.com/norwoodj/helm-docs/cmd/helm-docs

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -f helm-docs

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip-sign
