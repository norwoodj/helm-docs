# nginx-ingress

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

A simple wrapper around the stable/nginx-ingress chart that adds a few of our conventions

**Homepage:** <https://github.com/norwoodj/helm-docs/tree/master/example-charts/nginx-ingress>

## Maintainers

| Name | Email | URL |
| ---- | ------ | --- |
| John Norwood | <norwood.john.m@gmail.com> |  |

## Source Code

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/nginx-ingress>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| @stable | nginx-ingress | 0.22.1 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller | object | `{"image":{"repository":"nginx-ingress-controller","tag":"18.0831"},"name":"controller"}` | The controller |
| controller.image | object | `{"repository":"nginx-ingress-controller","tag":"18.0831"}` | The image of the controller |
| controller.image.repository | string | `"nginx-ingress-controller"` | The repository of the controller |
| controller.image.tag | string | `"18.0831"` | The tag of the image of the controller |
| controller.name | string | `"controller"` | The name of the controller |
| image | object | `{"pullPolicy":"IfNotPresent","registry":"docker.io","repository":"lucernae/django-sample","tag":"3.1"}` | Image map |
| image.pullPolicy | string | `"IfNotPresent"` | Image pullPolicy |
| image.registry | string | `"docker.io"` | Image registry |
| image.repository | string | `"lucernae/django-sample"` | Image repository |
| image.tag | string | `"3.1"` | Image tag |
| labels | map | map[client-name:my-boss project-name:awesome-project user/workload:true] | The deployment label |

