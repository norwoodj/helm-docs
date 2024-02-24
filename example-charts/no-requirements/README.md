# no-requirements

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square)

A simple chart that installs, let's say PrometheusRules, that needs no sub-charts

**Homepage:** <https://github.com/norwoodj/helm-docs/tree/master/example-charts/no-requirements>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| John Norwood | <norwood.john.m@gmail.com> |  |

## Source Code

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/no-requirements>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| rules.latency.percentiles.99.duration | string | `"5m"` | Duration for which the 99th percentile must be above the threshold to alert |
| rules.latency.percentiles.99.threshold | float | `1.5` | Threshold in seconds for our 99th percentile latency above which the alert will fire |
| rules.statusCodes.codes.5xx.duration | string | `"5m"` | Duration for which the percent of 5xx responses must be above the threshold to alert |
| rules.statusCodes.codes.5xx.threshold | float | `1.5` | Threshold percentage of 5xx responses above which the alert will fire |

