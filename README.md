helm-docs
=========
[![Go Report Card](https://goreportcard.com/badge/github.com/norwoodj/helm-docs)](https://goreportcard.com/report/github.com/norwoodj/helm-docs)

The helm-docs tool auto-generates documentation from helm charts into markdown files. The resulting
files contain metadata about their respective chart and a table with each of the chart's values, their defaults, and an
optional description parsed from comments.

The markdown generation is entirely [gotemplate](https://golang.org/pkg/text/template) driven. The tool parses metadata
from charts and generates a number of sub-templates that can be referenced in a template file (by default `README.md.gotmpl`).
If no template file is provided, the tool has a default internal template that will generate a reasonably formatted README.

The most useful aspect of this tool is the auto-detection of field descriptions from comments:
```yaml
config:
  databasesToCreate:
    # -- default database for storage of database metadata
    - postgres

    # -- database for the [hashbash](https://github.com/norwoodj/hashbash-backend-go) project
    - hashbash

  usersToCreate:
    # -- admin user
    - {name: root, admin: true}

    # -- user with access to the database with the same name
    - {name: hashbash, readwriteDatabases: [hashbash]}

statefulset:
  image:
    # -- Image to use for deploying, must support an entrypoint which creates users/databases from appropriate config files
    repository: jnorwood/postgresql
    tag: "11"

  # -- Additional volumes to be mounted into the database container
  extraVolumes:
    - name: data
      emptyDir: {}
```

Resulting in a resulting README section like so:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| config.databasesToCreate[0] | string | `"postgresql"` | default database for storage of database metadata |
| config.databasesToCreate[1] | string | `"hashbash"` | database for the [hashbash](https://github.com/norwoodj/hashbash-backend-go) project |
| config.usersToCreate[0] | object | `{"admin":true,"name":"root"}` | admin user |
| config.usersToCreate[1] | object | `{"name":"hashbash","readwriteDatabases":["hashbash"]}` | user with access to the database with the same name |
| statefulset.extraVolumes | list | `[{"emptyDir":{},"name":"data"}]` | Additional volumes to be mounted into the database container |
| statefulset.image.repository | string | `"jnorwood/postgresql:11"` | Image to use for deploying, must support an entrypoint which creates users/databases from appropriate config files |
| statefulset.image.tag | string | `"18.0831"` |  |

You'll notice that some complex fields (lists and objects) are documented while others aren't, and that some simple fields
like `statefulset.image.tag` are documented even without a description comment. The rules for what is and isn't documented in
the final table will be described in detail later in this document.

## Installation
helm-docs can be installed using [homebrew](https://brew.sh/):

```bash
brew install norwoodj/tap/helm-docs
```

This will download and install the [latest release](https://github.com/norwoodj/helm-docs/releases/latest)
of the tool.

To build from source in this repository:

```bash
cd cmd/helm-docs
go build
```

Or install from source:

```bash
GO111MODULE=on go get github.com/norwoodj/helm-docs/cmd/helm-docs
```

## Usage

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

## Ignoring Chart Directories
helm-docs supports a `.helmdocsignore` file, exactly like a `.gitignore` file in which one can specify directories to ignore
when searching for charts. Directories specified need not be charts themselves, so parent directories containing potentially
many charts can be ignored and none of the charts underneath them will be processed. You may also directly reference the
Chart.yaml file for a chart to skip processing for it.

## Markdown Rendering
There are two important paramaters to be aware of when running helm-docs. `--chart-serach-root` specifies the directory
under which the tool will recursively search for charts to render documentation for. `--template-files` specifies the list
of gotemplate files that should be used in rendering the resulting markdown file for each chart found. By default
`--chart-search-root=.` and `--template-files=README.md.gotmpl`.

If a template file is specified as a filename only as with the default above, the file is interpreted as being _relative to each chart directory found_.
If however a template file is specified as a relative path, e.g. the first of `--template-files=./_templates.gotmpl --template-files=README.md.gotmpl`
then the file is interpreted as being relative to the `chart-search-root`.

This repo is a good example of this in action. If you take a look at the [.pre-commit-config.yaml file](./.pre-commit-config.yaml)
here, you'll see our search root is set to [example-charts](./example-charts) and the list of templates used for each chart
is the [_templates.gotmpl file in that directory](./example-charts/_templates.gotmpl) and the README.md.gotmpl file in
each chart directory.

If any of the specified template files is not found for a chart (you'll notice most of the example charts do not have a README.md.gotmpl)
file, then the internal default template is used instead.

In addition to extra defined templates you specify in these template files, there are quite a few built-in templates that
can be used as well:

| Name | Description |
|------|-------------|
| chart.header              | The main heading of the generated markdown file |
| chart.name                | The _name_ field from the chart's `Chart.yaml` file |
| chart.deprecationWarning  | A deprecation warning which is displayed when the _deprecated_ field from the chart's `Chart.yaml` file is `true` |
| chart.description         | A description line containing the _description_ field from the chart's `Chart.yaml` file, or "" if that field is not set |
| chart.version             | The _version_ field from the chart's `Chart.yaml` file |
| chart.versionBadge        | A badge stating the current version of the chart |
| chart.type                | The _type_ field from the chart's `Chart.yaml` file |
| chart.typeBadge           | A badge stating the current type of the chart |
| chart.appVersion          | The _appVersion_ field from the chart's `Chart.yaml` file |
| chart.appVersionBadge     | A badge stating the current appVersion of the chart |
| chart.homepage            | The _home_ link from the chart's `Chart.yaml` file, or "" if that field is not set |
| chart.homepageLine        | A text line stating the current homepage of the chart |
| chart.maintainersHeader   | The heading for the chart maintainers section |
| chart.maintainersTable    | A table of the chart's maintainers |
| chart.maintainersSection  | A section headed by the maintainersHeader from above containing the maintainersTable from above or "" if there are no maintainers |
| chart.sourcesHeader       | The heading for the chart sources section |
| chart.sourcesList         | A list of the chart's sources |
| chart.sourcesSection      | A section headed by the sourcesHeader from above containing the sourcesList from above or "" if there are no sources |
| chart.kubeVersion         | The _kubeVersion_ field from the chart's `Chart.yaml` file |
| chart.kubeVersionLine     | A text line stating the required Kubernetes version for the chart |~~~~
| chart.requirementsHeader  | The heading for the chart requirements section |
| chart.requirementsTable   | A table of the chart's required sub-charts |
| chart.requirementsSection | A section headed by the requirementsHeader from above containing the kubeVersionLine and/or the requirementsTable from above or "" if there are no requirements |
| chart.valuesHeader        | The heading for the chart values section |
| chart.valuesTable         | A table of the chart's values parsed from the `values.yaml` file (see below) |
| chart.valuesSection       | A section headed by the valuesHeader from above containing the valuesTable from above or "" if there are no values |

The default internal template mentioned above uses many of these and looks like this:
```
{{ template "chart.header" . }}
{{ template "chart.deprecationWarning" . }}

{{ template "chart.badgesSection" . }}

{{ template "chart.description" . }}

{{ template "chart.homepageLine" . }}

{{ template "chart.maintainersSection" . }}

{{ template "chart.sourcesSection" . }}

{{ template "chart.requirementsSection" . }}

{{ template "chart.valuesSection" . }}
```

The tool also includes the [sprig templating library](https://github.com/Masterminds/sprig), so those functions can be used
in the templates you supply.

### values.yaml metadata
This tool can parse descriptions and defaults of values from `values.yaml` files. The defaults are pulled directly from
the yaml in the file. 

It was formerly the case that descriptions had to be specified with the full path of the yaml field. This is no longer
the case, although it is still supported. Where before you would document a values.yaml like so:

```yaml
controller:
  publishService:
    # controller.publishService.enabled -- Whether to expose the ingress controller to the public world
    enabled: false

  # controller.replicas -- Number of nginx-ingress pods to load balance between.
  # Do not set this below 2.
  replicas: 2
```

You may now equivelantly write:
```yaml
controller:
  publishService:
    # -- Whether to expose the ingress controller to the public world
    enabled: false

  # -- Number of nginx-ingress pods to load balance between.
  # Do not set this below 2.
  replicas: 2
```

New-style comments are much the same as the old-style comments, except that while old comments for a field could appear
anywhere in the file, new-style comments must appear **on the line(s) immediately preceding the field being documented.**

I invite you to check out the [example-charts](./example-charts) to see how this is done in practice. The `but-auto-comments`
examples in particular document the new comment format.

Note that comments can continue on the next line. In that case leave out the double dash, and the lines will simply be
appended with a space in-between, as in the `controller.replicas` field in the example above

The following rules are used to determine which values will be added to the values table in the README:

* By default, only _leaf nodes_, that is, fields of type `int`, `string`, `float`, `bool`, empty lists, and empty maps
  are added as rows in the values table. These fields will be added even if they do not have a description comment
* Lists and maps which contain elements will not be added as rows in the values table _unless_ they have a description
  comment which refers to them
* Adding a description comment for a non-empty list or map in this way makes it so that leaf nodes underneath the
  described field will _not_ be automatically added to the values table. In order to document both a non-empty list/map
  _and_ a leaf node within that field, description comments must be added for both

e.g. In this case, both `controller.livenessProbe` and `controller.livenessProbe.httpGet.path` will be added as rows in
the values table, but `controller.livenessProbe.httpGet.port` will not
```yaml
controller:
  # -- Configure the healthcheck for the ingress controller
  livenessProbe:
    httpGet:
      # -- This is the liveness check endpoint
      path: /healthz
      port: http
```

Results in:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.livenessProbe | object | `{"httpGet":{"path":"/healthz","port":8080}}` | Configure the healthcheck for the ingress controller |
| controller.livenessProbe.httpGet.path | string | `"/healthz"` | This is the liveness check endpoint |

If we remove the comment for `controller.livenessProbe` however, both leaf nodes `controller.livenessProbe.httpGet.path`
and `controller.livenessProbe.httpGet.port` will be added to the table, with or without description comments:

```yaml
controller:
  livenessProbe:
    httpGet:
      # -- This is the liveness check endpoint
      path: /healthz
      port: http
```

Results in:

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.livenessProbe.httpGet.path | string | `"/healthz"` | This is the liveness check endpoint |
| controller.livenessProbe.httpGet.port | string | `"http"` | |


### nil values
If you would like to define a key for a value, but leave the default empty, you can still specify a description for it
as well as a type. This is possible with both the old and the new comment format:
```yaml
controller:
  # -- (int) Number of nginx-ingress pods to load balance between
  replicas:
  
  # controller.image -- (string) Number of nginx-ingress pods to load balance between
  image:
```
This could be useful when wanting to enforce user-defined values for the chart, where there are no sensible defaults.

### Default values/column
In cases where you do not want to include the default value from `values.yaml`, or where the real default is calculated
inside the chart, you can change the contents of the column like so:

```yaml
service:
  # -- Add annotations to the service, this is going to be a long comment across multiple lines
  # but that's fine, these will be concatenated and the @default will be rendered as the default for this field
  # @default -- the chart will add some internal annotations automatically
  annotations: []
```

The order is important. The first comment line(s) must be the one specifying the key or using the auto-detection feature and
the description for the field. The `@default` comment must follow.

See [here](./example-charts/custom-template/values.yaml) for an example.

### Spaces and Dots in keys
In the old-style comment, if a key name contains any "." or " " characters, that section of the path must be quoted in
description comments e.g.

```yaml
service:
  annotations:
    # service.annotations."external-dns.alpha.kubernetes.io/hostname" -- Hostname to be assigned to the ELB for the service
    external-dns.alpha.kubernetes.io/hostname: stupidchess.jmn23.com

configMap:
  # configMap."not real config param" -- A completely fake config parameter for a useful example
  not real config param: value
```
