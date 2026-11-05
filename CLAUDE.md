# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

helm-docs is a Go CLI that generates `README.md` files for Helm charts from `Chart.yaml`, `values.yaml` (and its comments), and Go templates (`README.md.gotmpl`).

## Commands

```bash
make helm-docs                 # build ./helm-docs binary
make test                      # go test -v ./...
make lint                      # go fmt + gosec (-exclude G304)
make generate-example-charts   # rebuild binary and regenerate every example-charts/*/README.md

go test ./pkg/document -run TestSection -v           # single test
go test -run=^$ -bench=. ./cmd/helm-docs             # benchmarks (uses cmd/helm-docs/testdata/benchmark)
```

CI (`.github/workflows/build.yml`) runs gosec, a goreleaser snapshot build, and `scripts/check-example-charts-diff`, which runs `make generate-example-charts` and fails on any `git diff`. **Any change that alters generated output must commit the regenerated `example-charts/*/README.md` files.**

## Architecture

Pipeline, driven by `cmd/helm-docs/main.go`:

1. **Find charts** — `helm.FindChartDirectories` walks `--chart-search-root`, honoring `.helmdocsignore` (`pkg/util/ignore.go`).
2. **Parse each chart** (in parallel) — `helm.ParseChartInformation` (`pkg/helm/chart_info.go`) loads `Chart.yaml`, requirements/dependencies, the values file as a `yaml.Node` tree (stripping `# @ignore` nodes), and scans the raw values file line by line for **explicit** comments of the form `# key.path -- description`. Result: `ChartDocumentationInfo` with a `map[string]ChartValueDescription` keyed by dotted value path.
3. **Dependency values** (optional, `--document-dependency-values`) — `document.GetDependencyValues` pulls in subchart values found elsewhere in the search root, prefixed with the dependency alias/name.
4. **Build rows** — `pkg/document/values.go` recursively walks the `yaml.Node` tree producing `valueRow`s. Descriptions come from two sources: the explicit map from step 2, or **auto-comments** (`# -- description` directly above a key, read from `node.HeadComment` via `helm.ParseComment` in `pkg/helm/comment.go`). Explicit descriptions take precedence. Leaf values and documented objects/lists become rows; undocumented objects are descended into.
5. **Model + render** — `pkg/document/model.go` sorts rows (`--sort-values-order`), groups them by `@section`, and builds `chartTemplateData`. `pkg/document/template.go` defines all built-in `chart.*` named templates (e.g. `chart.valuesTable`, `chart.valuesSection`) and loads user template files; `generate.go` renders and writes the output file.

### Comment annotations

Parsed in `pkg/helm`: `# --` (description), `# @default --`, `# @raw`, `# @notationType --`, `# @section --`, `# @ignore`, and a `(type)` prefix in the description to override the type column. Regexes live at the top of `pkg/helm/chart_info.go`; the same syntax must work in both the explicit (line-scan) path and the auto-comment (`HeadComment`) path.

### Configuration

Flags are defined in `cmd/helm-docs/command_line.go` and bound to viper (env prefix `HELM_DOCS_`). Packages read config via **global `viper.GetString/GetBool`** (e.g. `values-file`, `output-file`, `sort-values-order`, `ignore-non-descriptions`), so tests in `pkg/` that depend on these must set them with `viper.Set`.

### Tests

- Unit tests build values from inline YAML strings (`parseYamlValues` in `pkg/document/values_test.go`) and assert on `valueRow` fields.
- `example-charts/` doubles as end-to-end golden tests — each directory exercises one feature (sections, custom notation types, umbrella dependencies, etc.). Add a chart there when adding a user-visible feature, then run `make generate-example-charts`.
