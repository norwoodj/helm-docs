# umbrella

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

A chart demonstrating that values documentation from child charts are aggregated on the parent chart.

## Requirements

| Repository | Name | Version |
|------------|------|---------|
|  | library | 0.1.0 |
|  | sub-a | 0.1.0 |
|  | sub-b | 0.1.0 |
|  | sub-c | 0.1.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| global | object | `{}` |  |
| global.myGlobalKey | string | `"my-global-value"` | A global key |
| global.myGlobalSubChartKey | string | `"my-global-sub-chart-value"` | A global key defined in a sub chart |
| myParentKey | string | `"my-parent-value"` | A parent key |
| sub-a.mySubKeyA | string | `"my-sub-value-a"` | Value for sub-chart A |
| sub-b.mySubKeyB | string | `"my-sub-value-b"` | Value for sub-chart B |

