# Sections

This creates values, but sectioned into own section tables if a section comment is provided.

## Values

### Some Section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | Uses the name of the configmap created by this chart |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. |
| controller.podLabels | object | `{}` | The labels to be applied to instances of the controller pod |

### Special Attention

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.ingressClass | string | `"nginx"` | You can also specify value comments like this |
| controller.publishService | object | <pre lang="yaml">enabled: false</pre> | This is a publishService |
| controller.replicas | int | `nil` | Number of nginx-ingress pods to load balance between |
| ingress.host | string | `""` | Overrides the inherited section |

### Inherited Section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| ingress.enabled | bool | `false` | Inherits "Inherited Section" from the parent key |

### Other Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].name | string | `"config-volume"` |  |
| controller.image.repository | string | `"nginx-ingress-controller"` |  |
| controller.image.tag | string | `"18.0831"` |  |
| controller.name | string | `"controller"` |  |
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |
| controller.service.type | string | `"LoadBalancer"` |  |

## Values

<h3>Some Section</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>controller.extraVolumes[0].configMap.name</td>
			<td>string</td>
			<td><pre lang="json">
"nginx-ingress-config"
</pre>
</td>
			<td>Uses the name of the configmap created by this chart</td>
		</tr>
		<tr>
			<td>controller.persistentVolumeClaims</td>
			<td>list</td>
			<td><pre lang="">
the chart will construct this list internally unless specified
</pre>
</td>
			<td>List of persistent volume claims to create.</td>
		</tr>
		<tr>
			<td>controller.podLabels</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>The labels to be applied to instances of the controller pod</td>
		</tr>
	</tbody>
</table>
<h3>Special Attention</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>controller.ingressClass</td>
			<td>string</td>
			<td><pre lang="json">
"nginx"
</pre>
</td>
			<td>You can also specify value comments like this</td>
		</tr>
		<tr>
			<td>controller.publishService</td>
			<td>object</td>
			<td><pre lang="yaml">
enabled: false
</pre>
</td>
			<td>This is a publishService</td>
		</tr>
		<tr>
			<td>controller.replicas</td>
			<td>int</td>
			<td><pre lang="json">
null
</pre>
</td>
			<td>Number of nginx-ingress pods to load balance between</td>
		</tr>
		<tr>
			<td>ingress.host</td>
			<td>string</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Overrides the inherited section</td>
		</tr>
	</tbody>
</table>
<h3>Inherited Section</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>ingress.enabled</td>
			<td>bool</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>Inherits "Inherited Section" from the parent key</td>
		</tr>
	</tbody>
</table>

<h3>Other Values</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
	<tr>
		<td>controller.extraVolumes[0].name</td>
		<td>string</td>
		<td><pre lang="json">
"config-volume"
</pre>
</td>
		<td></td>
	</tr>
	<tr>
		<td>controller.image.repository</td>
		<td>string</td>
		<td><pre lang="json">
"nginx-ingress-controller"
</pre>
</td>
		<td></td>
	</tr>
	<tr>
		<td>controller.image.tag</td>
		<td>string</td>
		<td><pre lang="json">
"18.0831"
</pre>
</td>
		<td></td>
	</tr>
	<tr>
		<td>controller.name</td>
		<td>string</td>
		<td><pre lang="json">
"controller"
</pre>
</td>
		<td></td>
	</tr>
	<tr>
		<td>controller.service.annotations."external-dns.alpha.kubernetes.io/hostname"</td>
		<td>string</td>
		<td><pre lang="json">
"stupidchess.jmn23.com"
</pre>
</td>
		<td>Hostname to be assigned to the ELB for the service</td>
	</tr>
	<tr>
		<td>controller.service.type</td>
		<td>string</td>
		<td><pre lang="json">
"LoadBalancer"
</pre>
</td>
		<td></td>
	</tr>
	</tbody>
</table>

