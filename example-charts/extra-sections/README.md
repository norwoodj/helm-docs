# extra-sections

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 13.0.0](https://img.shields.io/badge/AppVersion-13.0.0-informational?style=flat-square)

A chart for showing how to use extra sections

**Homepage:** <https://github.com/norwoodj/helm-docs/tree/master/example-charts/extra-sections>

## Source Code

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/extra-sections>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. See [Volume Claim Template](#volume-claim-template) for details. |
| controller.replicas | int | `1` | Number of controller replicas |
| controller.serviceAnnotations | object | `{}` | Annotations for the controller service. See [Service Annotations](#service-annotations) for examples. |

### Service Annotations
Annotations are passed as-is to the service. Some commonly used annotations:

| Annotation | Description |
|------------|-------------|
| `external-dns.alpha.kubernetes.io/hostname` | Hostname to register in DNS |
| `service.beta.kubernetes.io/aws-load-balancer-internal` | Use an internal AWS load balancer |

### Upgrading
Extra sections don't need to belong to a value. This one is defined on its own and
describes things to consider when upgrading the chart.

1. Read the release notes.
2. Check whether any of the values you set have been renamed.

### Volume Claim Template
The volume claim template is used to create a `PersistentVolumeClaim` for the controller
pods. The amount of storage requested is specified in the `.spec.resources.requests.storage`
field. The default value of `25Mi` should meet the requirements of most workloads.

You can use the following formula to roughly estimate the required storage:
```
( N_items * 200 ) + ( 64 * 1024 ) = total number of bytes
```

More information about persistent volume claims can be found in
[PersistentVolumeClaims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims).

