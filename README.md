helm-docs
=========
The helm-docs tool generates automatic documentation from a helm chart into a markdown file. The resulting
file contains metadata about the chart and a table with all of your chart's values, their defaults, and an
optional description parsed from comments.

To build:
```bash
$ make
```

The tool can be invoked on its own, but it must be done in a chart directory containing `values.yaml` and `Chart.yaml`
files.

```bash
$ helm-docs --dry-run # dry-run prints output to stdout rather than modifying the README in the directory you're in
```

## values.yaml metadata
This tool can parse descriptions and defaults of values from values.yaml. The defaults is done automatically by simply
parsing the yaml in the file. Descriptions can be added for parameters by specifying the full path of the value and
a particular format. I invite you to check out the example-chart to see how this is done in practice. In order to add
a description for a parameter you need only put a comment somewhere in the file of the format:

```yaml
controller:
  publishService:
    # controller.publishService.enabled -- Whether to expose the ingress controller to the public world
    enabled: false

  # controller.replicas -- Number of nginx-ingress pods to load balance between
  replicas: 2
```

And the descriptions will be picked up and put in the table in the README. The comment need not be near the parameter it
explains, although this is probably preferable.
