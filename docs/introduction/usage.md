# Usage
## How to run?

### Pre-commit hook

If you want to automatically generate `README.md` files with a pre-commit hook, make sure you
[install the pre-commit binary](https://pre-commit.com/#install), and add a [.pre-commit-config.yaml file](./.pre-commit-config.yaml)
to your project. Then run:

```bash
pre-commit install
pre-commit install-hooks
```

Future changes to your chart's `requirements.yaml`, `values.yaml`, `Chart.yaml`, or `README.md.gotmpl` files will cause an update to documentation when you commit.

### Running the binary directly

To run and generate documentation into READMEs for all helm charts within or recursively contained by a directory:

```bash
helm-docs
# OR
helm-docs --dry-run # prints generated documentation to stdout rather than modifying READMEs
```

The tool searches recursively through subdirectories of the current directory for `Chart.yaml` files and generates documentation
for every chart that it finds.

### Using docker

You can mount a directory with charts under `/helm-docs` within the container.

Then run:

```bash
docker run --rm --volume "$(pwd):/helm-docs" -u $(id -u) jnorwood/helm-docs:latest
```

## Specifics
### Ignoring Chart Directories
helm-docs supports a `.helmdocsignore` file, exactly like a `.gitignore` file in which one can specify directories to ignore
when searching for charts. Directories specified need not be charts themselves, so parent directories containing potentially
many charts can be ignored and none of the charts underneath them will be processed. You may also directly reference the
Chart.yaml file for a chart to skip processing for it.

### Generating Doc with Dependency values
Umbrella Helm chart documentation can include dependency values with `document-dependency-values` flag.
All dependency values will be merged into values of umbrella chart documentation.

If you want to include dependency values, but don't want to generate doc for each dependency:
* set `chart-search-root` parameter to directory that contains umbrella chart and all dependency charts.
* list all charts you want to generate doc using `chart-to-generate` flag
* set `document-dependency-values` flag to true
