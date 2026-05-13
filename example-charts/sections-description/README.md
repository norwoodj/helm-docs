# Sections

This creates values, but sectioned into own section tables if a section comment is provided.

Also demonstrates use of `@sectionDescription` annotation for sections, and propagation of section tag to children of lists and objects.

## Values

### Some Section
  
SomeSection has a description  
It will have multiple lines, and it will be rendered in the documentation with newlines preserved.  
  
It can have multiple paragraphs  
> It supports any markdown syntax
  
This is another paragraph for the same section  
It's better to use just one description tag per section, as the order depends on the sort, but you can have descriptions separated and related to different values if you want

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.extraVolumes[0].configMap.name | string | `"nginx-ingress-config"` | section gets propagated to this value,  to propagate, the ancestor needs to have `# --` comment even if it's not supposed to be added to the documentation  |
| controller.extraVolumes[0].name | string | `"config-volume"` |  |
| controller.persistentVolumeClaims | list | the chart will construct this list internally unless specified | List of persistent volume claims to create. The Regular description continues here |
| controller.podLabels | object | empty map | The labels to be applied to instances of the controller pod |
| controller.service.annotations."external-dns.alpha.kubernetes.io/hostname" | string | `"stupidchess.jmn23.com"` | Hostname to be assigned to the ELB for the service |

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

### Different nested Section

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.service.type | string | `"LoadBalancer"` | Nested section can be different if specified |

### Other Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| controller.image.repository | string | `"nginx-ingress-controller"` |  |
| controller.image.tag | string | `"18.0831"` |  |
| controller.name | string | `"controller"` |  |

## Values

<h3>Some Section</h3>
  
SomeSection has a description  
It will have multiple lines, and it will be rendered in the documentation with newlines preserved.  
  
It can have multiple paragraphs  
> It supports any markdown syntax
  
This is another paragraph for the same section  
It's better to use just one description tag per section, as the order depends on the sort, but you can have descriptions separated and related to different values if you want

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
			<td>section gets propagated to this value,  to propagate, the ancestor needs to have `# --` comment even if it's not supposed to be added to the documentation </td>
		</tr>
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
			<td>controller.persistentVolumeClaims</td>
			<td>list</td>
			<td><pre lang="">
the chart will construct this list internally unless specified
</pre>
</td>
			<td>List of persistent volume claims to create. The Regular description continues here</td>
		</tr>
		<tr>
			<td>controller.podLabels</td>
			<td>object</td>
			<td><pre lang="">
empty map
</pre>
</td>
			<td>The labels to be applied to instances of the controller pod</td>
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
	</tbody>
</table>
<h3>Special Attention</h3>
  
This section has a list:  
* Item 1  
* Item 2  
* Item 3

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
			<td><pre lang="json">
{
  "enabled": false
}
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
	</tbody>
</table>
<h3>Different nested Section</h3>
<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>controller.service.type</td>
			<td>string</td>
			<td><pre lang="json">
"LoadBalancer"
</pre>
</td>
			<td>Nested section can be different if specified</td>
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
	</tbody>
</table>

