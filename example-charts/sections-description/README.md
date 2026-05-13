# Sections

This creates values, but sectioned into own section tables if a section comment is provided.

## Values

### Some Section

  
SomeSection has a description  
It will have multiple lines, and it will be rendered in the documentation with newlines preserved.  
  
It can have multiple paragraphs - and can be closed with a closing tag, so you can continue with regular description
  
This is another paragraph for the same section - this one is closed by a different tag  
It's better to use just one description tag per section, as the order depends on the sort, but you can have descriptions separated and related to different values if you want

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | Uses the name of the configmap created by this chart |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. The Regular description continues here |
| controller.podLabels | object | empty map | The labels to be applied to instances of the controller pod |

### Special Attention

  
This section has a list:  
* Item 1  
* Item 2  
* Item 3

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.ingressClass | string | `"nginx"` | You can also specify value comments like this |
| controller.publishService | object | `{"enabled":false}` | This is a publishService |
| controller.replicas | int | `nil` | Number of nginx-ingress pods to load balance between |

### Other Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].name | string | `"config-volume"` |  |
| controller.image.repository | string | `"nginx-ingress-controller"` |  |
| controller.image.tag | string | `"18.0831"` |  |
| controller.name | string | `"controller"` |  |
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |
| controller.service.type | string | `"LoadBalancer"` |  |
