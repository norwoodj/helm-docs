helm-docs:
	go build github.com/norwoodj/helm-docs/cmd/helm-docs

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

.PHONY: clean
clean:
	rm -f helm-docs

.PHONY: dist
dist:
	goreleaser release --rm-dist --snapshot --skip=sign
