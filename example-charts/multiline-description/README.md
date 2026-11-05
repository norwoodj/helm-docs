# multiline-description

Multiline @raw descriptions with code snippets rendered in the markdown values table

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| image | string | `"nginx"` | A regular one-line description |
| password | string | a random 32 character string | Postgresql user password.<br>Generated when installed for the first time. |
| podAnnotations | object | `{}` | Extra annotations for the pod<br>Every comment line is a new line in the table.<br>Code fences become pre blocks, so pipes and angle brackets inside them are safe:<br><pre lang="yaml">podAnnotations:<br>  example.com/script: \|<br>    echo &#34;&lt;hello&gt;&#34;</pre> |
| replicas | int | `1` | Number of replicas<br>Explicit key comments work too.<br><pre lang="bash">kubectl scale deploy/app --replicas=3</pre> |

The same values rendered as an HTML table:

## Values

<table>
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td>image</td>
			<td>string</td>
			<td><pre lang="json">
"nginx"
</pre>
</td>
			<td>A regular one-line description</td>
		</tr>
		<tr>
			<td>password</td>
			<td>string</td>
			<td><pre lang="">
a random 32 character string
</pre>
</td>
			<td>Postgresql user password.<br>Generated when installed for the first time.</td>
		</tr>
		<tr>
			<td>podAnnotations</td>
			<td>object</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Extra annotations for the pod<br>Every comment line is a new line in the table.<br>Code fences become pre blocks, so pipes and angle brackets inside them are safe:<br><pre lang="yaml">podAnnotations:<br>  example.com/script: |<br>    echo &#34;&lt;hello&gt;&#34;</pre></td>
		</tr>
		<tr>
			<td>replicas</td>
			<td>int</td>
			<td><pre lang="json">
1
</pre>
</td>
			<td>Number of replicas<br>Explicit key comments work too.<br><pre lang="bash">kubectl scale deploy/app --replicas=3</pre></td>
		</tr>
	</tbody>
</table>

