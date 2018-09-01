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
a particular format:

```hcl-terraform
controller:
  name: controller
  image:
    repository: nginx-ingress-controller
    tag: "18.0831"

  # controller.ingressClass -- Name of the ingress class to route through this controller
  ingressClass: nginx

  # controller.podLabels -- The labels to be applied to instances of the controller pod
  podLabels: {}

  publishService:
    # controller.publishService.enabled -- Whether to expose the ingress controller to the public world
    enabled: false
```

This would produce the following table:
```
| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.name | string | "controller" |  |
| controller.image.repository | string | "nginx-ingress-controller" |  |
| controller.image.tag | string | "18.0831" |  |
| controller.ingressClass | string | "nginx" | Name of the ingress class to route through this controller |
| controller.podLabels | object | {} | The labels to be applied to instances of the controller pod |
| ingresses.applicationBaseUrl | string | "sre.rmneng.com" | The base URL to use
