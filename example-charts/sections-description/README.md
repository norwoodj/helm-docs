# Sections

This creates values, but sectioned into own section tables if a section comment is provided.

## Values

### Some Section
Some Section description

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | Uses the name of the configmap created by this chart |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. |
| controller.podLabels | object | `{}` | The labels to be applied to instances of the controller pod |

### Special Attention
Special Attention description

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.ingressClass | string | `"nginx"` | You can also specify value comments like this |
| controller.publishService | object | `{"enabled":false}` | This is a publishService |
| controller.replicas | int | `nil` | Number of nginx-ingress pods to load balance between |

### Other Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |
