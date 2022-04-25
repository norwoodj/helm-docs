# full-template

## `extra.flower`
```
          ,-.
 ,     ,-.   ,-.
/ \   (   )-(   )
\ |  ,.>-(   )-<
 \|,' (   )-(   )
  Y ___`-'   `-'
  |/__/   `-'
  |
  |
  |    -hi-
__|_____________
```

## `chart.deprecationWarning`
> **:exclamation: This Helm Chart is deprecated!**

## `chart.name`

full-template

## `chart.description`

A chart for showing every README-element

## `chart.version`

1.0.0

## `chart.versionBadge`

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square)

## `chart.type`

application

## `chart.typeBadge`

![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

## `chart.appVersion`

13.0.0

## `chart.appVersionBadge`

![AppVersion: 13.0.0](https://img.shields.io/badge/AppVersion-13.0.0-informational?style=flat-square)

## `chart.badgesSection`

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 13.0.0](https://img.shields.io/badge/AppVersion-13.0.0-informational?style=flat-square)

## `chart.homepage`

https://github.com/norwoodj/helm-docs/tree/master/example-charts/full-template

## `chart.homepageLine`

**Homepage:** <https://github.com/norwoodj/helm-docs/tree/master/example-charts/full-template>

## `chart.maintainersHeader`

## Maintainers

## `chart.maintainersTable`

| Name | Email | Url |
| ---- | ------ | --- |
| John Norwood | <norwood.john.m@gmail.com> |  |

## `chart.maintainersSection`

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| John Norwood | <norwood.john.m@gmail.com> |  |

## `chart.sourcesHeader`

## Source Code

## `chart.sourcesList`

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/full-template>

## `chart.sourcesSection`

## Source Code

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/full-template>

## `chart.kubeVersion`

<=1.18

## `chart.kubeVersionLine`

Kubernetes: `<=1.18`

## `chart.requirementsHeader`

## Requirements

## `chart.requirementsTable`

| Repository | Name | Version |
|------------|------|---------|
| @stable | nginx-ingress | 0.22.1 |

## `chart.requirementsSection`

## Requirements

Kubernetes: `<=1.18`

| Repository | Name | Version |
|------------|------|---------|
| @stable | nginx-ingress | 0.22.1 |

## `chart.valuesHeader`

## Values

## `chart.valuesTable`

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | Uses the name of the configmap created by this chart |
| controller.extraVolumes[0].name | string | `"config-volume"` |  |
| controller.image.repository | string | `"nginx-ingress-controller"` |  |
| controller.image.tag | string | `"18.0831"` |  |
| controller.ingressClass | string | `"nginx"` | Name of the ingress class to route through this controller |
| controller.name | string | `"controller"` |  |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. For very long comments, break them into multiple lines. |
| controller.podLabels | object | `{}` | The labels to be applied to instances of the controller pod |
| controller.publishService.enabled | bool | `false` | Whether to expose the ingress controller to the public world |
| controller.replicas | int | `nil` | Number of nginx-ingress pods to load balance between |
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |
| controller.service.type | string | `"LoadBalancer"` |  |

## `chart.valuesSection`

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | Uses the name of the configmap created by this chart |
| controller.extraVolumes[0].name | string | `"config-volume"` |  |
| controller.image.repository | string | `"nginx-ingress-controller"` |  |
| controller.image.tag | string | `"18.0831"` |  |
| controller.ingressClass | string | `"nginx"` | Name of the ingress class to route through this controller |
| controller.name | string | `"controller"` |  |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. For very long comments, break them into multiple lines. |
| controller.podLabels | object | `{}` | The labels to be applied to instances of the controller pod |
| controller.publishService.enabled | bool | `false` | Whether to expose the ingress controller to the public world |
| controller.replicas | int | `nil` | Number of nginx-ingress pods to load balance between |
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |
| controller.service.type | string | `"LoadBalancer"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.9.1](https://github.com/norwoodj/helm-docs/releases/v1.9.1)
