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

or [scoop](https://scoop.sh):

```bash
scoop install helm-docs
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

## Generating Doc with Dependency values
Umbrella Helm chart documentation can include dependency values with `document-dependency-values` flag.
All dependency values will be merged into values of umbrella chart documentation.

If you want to include dependency values, but don't want to generate doc for each dependency:
* set `chart-search-root` parameter to directory that contains umbrella chart and all dependency charts.
* list all charts you want to generate doc using `chart-to-generate` flag
* set `document-dependency-values` flag to true

## Markdown Rendering
There are two important parameters to be aware of when running helm-docs. `--chart-search-root` specifies the directory
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
| chart.valuesTableHtml     | Like `chart.valuesTable` but it is rendered as (X)HTML tags to allow further rendering customization, instead of markdown tables format. |
| chart.valuesSectionHtml   | Like `chart.valuesSection` but uses `chart.valuesTableHtml` |
| chart.valueDefaultColumnRender | This is a hook template if you want to redefine how helm-docs render the default values in `chart.valuesTableHtml` mode. This is especially useful when combined with (X)HTML tags, so that you can nicely format multiline default values, like YAML/JSON object tree snippet with codeblock syntax highlighter, which is not possible or difficult when using the markdown table format. It can be redefined in your template file. |

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

You may now equivalently write:
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
### Ignoring values
In cases you would like to ignore certain values, you can mark it with @ignored tag:

```yaml
# @ignored
service:
  port: 8080
```

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

### Advanced table rendering
Some helm chart `values.yaml` uses complicated structure for the key/value
pairs. For example, it may uses a multiline string of Go template text instead
of plain strings. Some values might also refer to a certain YAML/JSON object
structure, like internal k8s value type, or an enum. For these use case,
a standard markdown table format might be inadequate and you want to use HTML
tags to render the table.

Some example use case on why you need advanced table rendering:

 - Hyperlinking the value type to an anchor or HTML link somewhere for reference
 - Collapsible value description using `<summary>` tags to save space
 - Multiline default values as codeblocks, instead of one line JSON structure for readability
 - Custom rendering, for colors, actions, bookmarking, cross-reference, etc
 - Cascading the markdown file generated by helm-docs to be post-processed by Jamstack into a static HTML docs site.

In order to accomodate this, `helm-docs` provides an extensible and flexible way to customize rendering.

1. Use the HTML value renderer instead of the default markdown format

You can use `chart.valuesSectionHtml` to render the values table as HTML tags,
instead of using `chart.valuesSection`. Using HTML tables provides more
flexibility because it can be processed by markdown viewer as a nested blocks,
instead of one row per line. This allows you to customize how each columns in a
row are rendered.

2. Overriding built-in templates

You can always overrides or redefine built-in templates in your own `_templates.
gotmpl` file. The built-in templates can be thought of as a template hook.
For example, if you need to change the HTML table, for example to add a new
column, or define maximum width/height, you can override `chart.valuesTableHtml`. Your overrides will then be called by `chart.valuesSectionHtml`.

You can add your own rendering logic for each column. For example, we have `chart.valueDefaultColumnRender` that is used to render "default value" column for each rows. If you want to override how helm-docs render the
"type" column, just define your own rendering template and call it from
`chart.valuesTableHtml` for each of the rows.

3. Using the metadata of each rows of values

Custom styling and rendering can be done as flexible as you want, but you
still need a metadata that describes each rows of values. You can access
this information from the templates.

When you override `chart.valuesTableHtml`, as you can see in the original
definition in `func getValuesTableTemplates()` [pkg/document/template.go](pkg/document/template.go), we iterates each row of values.
For each "Value", it is modeled as a struct defined in `valueRow` struct
in [pkg/document/model.go](pkg/document/model.go). You can then use the
fields in your template.

Some fields here are directly referenced from `values.yaml`:
- `Key`: the full name of the key referenced in `values.yaml`
- `Type`: the type of the value of the key in `values.yaml`. Can be automatically inferred from YAML structure, or annotated using `# -- (mytype)` where `mytype` can be any string that you refer as the type of the value.
- `NotationType`: the notation of the type used to render the default value. If `Type` refers to the data type of the value, then `NotationType` refers to **how** this value should be written/rendered by helm-docs. Generally helm-docs only remembers the notation type, but it was the writer's responsibility to make a template tag to render a specific notation type. Annotate the key with `# @notationType -- (mynotation)` where `mynotation` is an identifier to tell the renderer how to write the value.
- `Default`: this is the default value of the key, found from `values.yaml`. It is either inferred from the YAML structure or defined using `# @default -- my default value` annotation, in case you need to show other example values.
- `Description`: this is the description of the key/value, taken from the comments found in the `values.yaml` for the referred key.
- `LineNumber`: this is the line number associated with where the key is declared. You can use this to construct an anchor to the actual `values.yaml` file.

Note that helm-docs only provides these information, but the default behaviour is to always render it in plain Markdown file to be viewed locally.

4. Use markdown files generated by helm-docs as intermediary files to be processed further

Public helm charts sometimes needs to be published as static content
instead of just stored in a repository. This is needed for helm users to
be able to view or browse the chart options and dependencies.

It is often more than enough to just browse the chart values options on
git hosting that is able to render markdown files as a nice HTML page, like GitHub or GitLab.
However, for a certain use case, you may want to use your own
documentation generator to host or publish the output of helm-docs.

If you use some kind of Jamstack like Gatsby or Hugo, you can use the
output of helm-docs as an input for these doc generator. A typical use
case is to override helm-docs built-in template so that it renders a
markdown or markdownX files to be processed by Gatsby or Hugo into
a static Web/Javascript page.

For a more concrete examples on how to do these custom rendering, see [example here](./example-charts/custom-value-notation-type/README.md)
