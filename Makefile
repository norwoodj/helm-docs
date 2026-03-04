.PHONY: helm-docs
helm-docs:
	go build github.com/norwoodj/helm-docs/cmd/helm-docs

.PHONY: install
install:
	go install github.com/norwoodj/helm-docs/cmd/helm-docs

.PHONY: generate-example-charts
generate-example-charts: helm-docs
	./helm-docs --chart-search-root=example-charts --template-files=./_templates.gotmpl --template-files=README.md.gotmpl --document-dependency-values

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: gosec
gosec:
	@which gosec > /dev/null || go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude G304 ./...

.PHONY: lint
lint: fmt gosec

.PHONY: clean
clean:
	rm -f helm-docs

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip=sign

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  helm-docs    - Build the helm-docs binary"
	@echo "  install      - Install the helm-docs binary"
	@echo "  fmt          - Format Go code"
	@echo "  test         - Run tests"
	@echo "  gosec        - Run security scan with gosec"
	@echo "  lint         - Run all linters (fmt, gosec)"
	@echo "  clean        - Clean build artifacts"
	@echo "  dist         - Create distribution with goreleaser"
	@echo "  help         - Show this help message"
