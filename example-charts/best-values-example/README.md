# best-values-example

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)

One of the best values parsing example charts here, exhibits several more complicated examples

**Homepage:** <https://github.com/norwoodj/helm-docs/tree/master/example-charts/best-values-example>

## Maintainers

| Name | Email | URL |
| ---- | ------ | --- |
| John Norwood | <norwood.john.m@gmail.com> |  |

## Source Code

* <https://github.com/norwoodj/helm-docs/tree/master/example-charts/best-values-example>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| config.databasesToCreate[0] | string | `"postgresql"` | default database for storage of database metadata |
| config.databasesToCreate[1] | string | `"hashbash"` | database for the [hashbash](https://github.com/norwoodj/hashbash) project |
| config.usersToCreate[0] | object | <pre lang="yaml">admin: true<br>name: root</pre> | admin user |
| config.usersToCreate[1] | object | <pre lang="yaml">name: hashbash<br>readwriteDatabases:<br>  - hashbash</pre> | user with access to the database with the same name |
| statefulset.extraVolumes | list | <pre lang="yaml">- emptyDir: {}<br>  name: data</pre> | Additional volumes to be mounted into the database container |
| statefulset.image.repository | string | `"jnorwood/postgresq"` | Image to use for deploying, must support an entrypoint which creates users/databases from appropriate config files |
| statefulset.image.tag | string | `"11"` |  |
| statefulset.livenessProbe | object | <pre lang="yaml">enabled: false</pre> | Configure the healthcheck for the database |
| statefulset.podLabels | object | `{}` | The labels to be applied to instances of the database |

