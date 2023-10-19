## Strict linting

Sometimes you might want to enforce helm-docs to fail when some values are not documented correctly.

By default, this option is turned off:

```shell
./helm-docs -c  example-charts/helm-3
INFO[2023-06-29T07:54:29-07:00] Found Chart directories [.]                  
INFO[2023-06-29T07:54:29-07:00] Generating README Documentation for chart example-charts/helm-3 
```

but you can use the `-x` flag to turn it on:

```shell
helm-docs -x -c  example-charts/helm-3
INFO[2023-06-29T07:55:12-07:00] Found Chart directories [.]                  
WARN[2023-06-29T07:55:12-07:00] Error parsing information for chart ., skipping: values without documentation: 
controller
controller.name
controller.image
controller.extraVolumes.[0].name
controller.extraVolumes.[0].configMap
controller.extraVolumes.[0].configMap.name
controller.livenessProbe.httpGet
controller.livenessProbe.httpGet.port
controller.publishService
controller.service
controller.service.annotations
controller.service.annotations.external-dns.alpha.kubernetes.io/hostname
```

The CLI also supports excluding fields by regexp using the `-z` argument

```shell
helm-docs -x -z="controller.*" -c  example-charts/helm-3
INFO[2023-06-29T08:18:55-07:00] Found Chart directories [.]                  
INFO[2023-06-29T08:18:55-07:00] Generating README Documentation for chart example-charts/helm-3 
```

Multiple regexp can be passed, as in the following example:

```shell
helm-docs -x -z="controller.image.*" -z="controller.service.*"  -z="controller.extraVolumes.*"  -c  example-charts/helm-3
INFO[2023-06-29T08:21:04-07:00] Found Chart directories [.]                  
WARN[2023-06-29T08:21:04-07:00] Error parsing information for chart ., skipping: values without documentation: 
controller
controller.name
controller.livenessProbe.httpGet
controller.livenessProbe.httpGet.port
controller.publishService 
```

It is also possible to ignore specific errors using the `-y` argument.

```shell
helm-docs -x -y="controller.name" -y="controller.service"  -c  example-charts/helm-3
INFO[2023-06-29T08:23:40-07:00] Found Chart directories [.]                  
WARN[2023-06-29T08:23:40-07:00] Error parsing information for chart ., skipping: values without documentation: 
controller
controller.image
controller.extraVolumes.[0].name
controller.extraVolumes.[0].configMap
controller.extraVolumes.[0].configMap.name
controller.livenessProbe.httpGet
controller.livenessProbe.httpGet.port
controller.publishService
controller.service.annotations
controller.service.annotations.external-dns.alpha.kubernetes.io/hostname
 
```
