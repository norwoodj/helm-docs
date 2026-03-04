# no-values

A very simple chart that doesn't even need any values for customization

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square)

## Additional Information

### Snapshot classes

| name     | resource group | incremental |
|----------|----------------|-------------|
| azure-incr | true |  |
| azure | false |  |
| local-incr | true |  |
| local | false |  |
| mask-data-incr | true | rg-mask-data |
| mask-data | false | rg-mask-data |

### Default resources

```yaml
requests:
    cpu: 10m
    memory: 100m
```

## Installing the Chart

To install the chart with the release name `my-release`:

```console
$ helm repo add foo-bar http://charts.foo-bar.com
$ helm install my-release foo-bar/no-values
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| resources.requests.cpu | string | `"10m"` |  |
| resources.requests.memory | string | `"100m"` |  |
| volumeSnapshotClass.azure.parameters | object | `{}` |  |
| volumeSnapshotClass.local.parameters | object | `{}` |  |
| volumeSnapshotClass.mask-data.parameters.resourceGroup | string | `"rg-mask-data"` |  |

