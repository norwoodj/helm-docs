# django

![Version: 0.2.1](https://img.shields.io/badge/Version-0.2.1-informational?style=flat-square) ![AppVersion: 3.1](https://img.shields.io/badge/AppVersion-3.1-informational?style=flat-square)

Generic chart for basic Django-based web app

**Homepage:** <https://www.djangoproject.com/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Rizky Maulana Nugraha | <lana.pcfre@gmail.com> |  |

## Source Code

* <https://github.com/django/django>

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| file://../../common/v1.0.0 | common | 1.0.0 |
| file://../../postgis/v0.2.1 | postgis | 0.2.1 |

# Some Long Description

This is a sample README with custom overrides.
Check the template in [README.md.gotmpl](README.md.gotmpl).

In that file, we redefine the template definition of `chart.valueDefaultColumnRender`
for some custom `@notationType` such as `string/email`.

This chart README uses `chart.valuesSectionHtml` instead of `chart.valuesSection`.
Using HTML table directly instead of using Markdown table allows us to control the table
presentation, such as the height. This is especially useful for very long `values.yaml` file,
and you need to scroll both horizontally and vertically to navigate the values.

In the template file, we redefine `chart.valuesTableHtml` so that we use table height of
400px at most. Github can understand that attribute. The more sophisticated use case is if you
want to combine helm-docs with a Jamstack static generator where you can have your own page generated
from this README.

The customization can goes even further. Normally, you can't define anchor in markdown unless it is a heading. But you can do so easily using HTML tags.
You can override the column key renderer by adding an `id` attribute so that it can be referred.
This way, when you write markdown links like [ingress.tls.secretName](#ingress--tls--secretName), clicking the link
will take you to the value description row.

## Value Types

One of the benefit of using HTML table is we can make a simple tooltip and anchor.
For example, the value [global.adminEmail](#global--adminEmail) is annotated as type `string/email`. We create
the definition of the value type here and can be anchored by links with `#stringemail` hyperlinks.

We can also create custom type column renderer, where we can assign a tooltip for each type.
Try this out. Go navigate to [global.adminEmail](#global--adminEmail) value, hover on the value type `string/email`, you will then see
some tooltip. Clicking the type link will direct you back to it's relevant value type section below.

Other useful case is If the type is a known type, like
Kubernetes service type, you can anchor the type to redirect user to k8s documentation page to learn more.
Check the value [persistence.staticDir.accessModes](#persistence--staticDir--accessModes)

### string/email

This value type is for a valid email address format. Such as owner@somedomain.org.

## Notation Type

Another reason to use HTML table is because in some cases we want to custom-render the default value.

In helm chart templates, sometimes author designs the template to accept a go template string value.
That means, the template string can be processed by helm chart and be replaced with dynamic computed values, before it was
rendered to the chart. Although it is very useful and flexible to make the default value be dynamic,
it is not entirely obvious for the chart users to see a go template as value in a `values.yml`.
It would then be helpful to custom-render these default values in the helm README, so that it is not
treated as a pure JSON object (because the syntax highlighter would be incorrect).
Instead we can custom render the presentation so it would make sense to the user.

In our example here, any key with a type `tpl/xxx` would be rendered as `<pre></pre>`
HTML tag, in which we both put the key name and the YAML multiline modifier `|` to make
it really clear that the key accept a multiline string as value, because it would be rendered as
YAML object by helm after the values are interpolated/substituted.

Take a look at [extraPodEnv](#extraPodEnv). The `Default` column shows the key name `extraPodEnv`, the multiline YAML
modifier `|`, and the template string which contains some go string template syntax `{{ }}`.

You can also control the HTML styling directly. In some markdown viewer, the HTML tag and inline styles
are respected, so the custom styles can be seen. Combined with a Jamstack approach, you can
design your template to also incorporate some custom React styles or simple CSS.

In our example here, [global.adminEmail](#global--adminEmail) is annotated with `email` notationType.
This allows you to insert custom rendering code for email. For supported markdown viewer, like Visual Studio Code,
the default value will have `green` color, and if clicked will direct you to your default email composer.

The reason we have two separate annotation, value type and notation type, is because several different types
can have the same type renderer. For example, any type `tpl/xxx` is a go template string, so it will be rendered the same
in our docs if we annotate it with `@notationType -- tpl`.

## Customized Rendering

This README also shows some possible customization with helm-docs. In the [README.md.gotmpl](README.md.gotmpl)
file, you can see that we modified the column `Key` to also be hyperlinked with the definition in `values.yaml`.
If you view this README.md files in GitHub and click the value's key, you will be directed to the
key location in the `values.yaml` file.

You can also render a raw string into the comments using `@raw` annotations.
You can jump to [sampleYaml](#sampleYaml) key and check it's description where it
uses HTML `<summary>` tag to collapse some part of the comments.

## Values

<table height="400px" >
	<thead>
		<th>Key</th>
		<th>Type</th>
		<th>Default</th>
		<th>Description</th>
	</thead>
	<tbody>
		<tr>
			<td id="extraConfigMap"><a href="./values.yaml#L111">extraConfigMap</a></td>
			<td>
tpl/dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraConfigMap: |
 
</pre>
</div>
			</td>
			<td>Define this for extra config map to be included in django-shared-config</td>
		</tr>
		<tr>
			<td id="extraPodEnv"><a href="./values.yaml#L88">extraPodEnv</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraPodEnv: |
  - name: DJANGO_SETTINGS_MODULE
    value: "django.settings"
  - name: DEBUG
    value: {{ .Values.global.debug | quote }}
  - name: ROOT_URLCONF
    value: {{ .Values.global.rootURLConf | quote }}
  - name: MAIN_APP_NAME
    value: {{ .Values.global.mainAppName | quote }}
 
</pre>
</div>
			</td>
			<td>Define this for extra Django environment variables</td>
		</tr>
		<tr>
			<td id="extraPodSpec"><a href="./values.yaml#L100">extraPodSpec</a></td>
			<td>
tpl/object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraPodSpec: |
 
</pre>
</div>
			</td>
			<td>This will be evaluated as pod spec</td>
		</tr>
		<tr>
			<td id="extraSecret"><a href="./values.yaml#L106">extraSecret</a></td>
			<td>
tpl/dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraSecret: |
 
</pre>
</div>
			</td>
			<td>Define this for extra secrets to be included in django-shared-secret secret</td>
		</tr>
		<tr>
			<td id="extraVolume"><a href="./values.yaml#L125">extraVolume</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraVolume: |
 
</pre>
</div>
			</td>
			<td>Define this for extra volume (in pair with extraVolumeMounts)</td>
		</tr>
		<tr>
			<td id="extraVolumeMounts"><a href="./values.yaml#L116">extraVolumeMounts</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
extraVolumeMounts: |
 
</pre>
</div>
			</td>
			<td>Define this for extra volume mounts in the pod</td>
		</tr>
		<tr>
			<td id="global"><a href="./values.yaml#L14">global</a></td>
			<td>
object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{
  "adminEmail": "admin@localhost",
  "adminPassword": {
    "value": null,
    "valueFrom": {
      "secretKeyRef": {
        "key": "admin-password",
        "name": null
      }
    }
  },
  "adminUser": "admin",
  "databaseHost": "postgis",
  "databaseName": "django",
  "databasePassword": {
    "value": null,
    "valueFrom": {
      "secretKeyRef": {
        "key": "database-password",
        "name": null
      }
    }
  },
  "databasePort": 5432,
  "databaseUsername": "django_db_user",
  "debug": "False",
  "djangoArgs": "[\"uwsgi\",\"--chdir=${REPO_ROOT}\",\"--module=${MAIN_APP_NAME}.wsgi\",\"--socket=:8000\",\"--http=0.0.0.0:8080\",\"--processes=5\",\"--buffer-size=8192\"]\n",
  "djangoCommand": "[\"/opt/django/scripts/docker-entrypoint.sh\"]\n",
  "djangoSecretKey": {
    "value": null,
    "valueFrom": {
      "secretKeyRef": {
        "key": "django-secret",
        "name": null
      }
    }
  },
  "djangoSettingsModule": "django.settings",
  "existingSecret": "",
  "mainAppName": "django",
  "mediaRoot": "/opt/django/media",
  "nameOverride": "django",
  "rootURLConf": "django.urls",
  "sharedSecretName": "django-shared-secret",
  "siteName": "django",
  "staticRoot": "/opt/django/static"
}
</pre>
</div>
			</td>
			<td>This key name is used for service interconnection between subcharts and parent charts.</td>
		</tr>
		<tr>
			<td id="global--adminEmail"><a href="./values.yaml#L43">global.adminEmail</a></td>
			<td>
<a href="#stringemail" title="
This value type is for a valid email address format. Such as owner@somedomain.org.">string/email</a>
</td>
			<td>
				<div style="max-width: 300px;">
<a href="mailto:admin@localhost" style="color: green;">"admin@localhost"</a>
</div>
			</td>
			<td>Default admin email sender</td>
		</tr>
		<tr>
			<td id="global--adminPassword--value"><a href="./values.yaml#L36">global.adminPassword.value</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
null
</pre>
</div>
			</td>
			<td>Specify this password value. If not, it will be autogenerated everytime chart upgraded</td>
		</tr>
		<tr>
			<td id="global--adminUser"><a href="./values.yaml#L33">global.adminUser</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"admin"
</pre>
</div>
			</td>
			<td>Default super user admin username</td>
		</tr>
		<tr>
			<td id="global--databaseHost"><a href="./values.yaml#L63">global.databaseHost</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"postgis"
</pre>
</div>
			</td>
			<td>Django database host location. By default this chart can generate standard postgres chart. So you can leave it as default. If you use external backend,  you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databaseName"><a href="./values.yaml#L61">global.databaseName</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django"
</pre>
</div>
			</td>
			<td>Django database name</td>
		</tr>
		<tr>
			<td id="global--databasePassword--value"><a href="./values.yaml#L55">global.databasePassword.value</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
null
</pre>
</div>
			</td>
			<td>Specify this password value. If not, it will be autogenerated everytime chart upgraded. If you use external backend, you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databasePort"><a href="./values.yaml#L65">global.databasePort</a></td>
			<td>
int
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
5432
</pre>
</div>
			</td>
			<td>Django database port. By default this chart can generate standard postgres chart. So you can leave it as default. If you use external backend,  you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databaseUsername"><a href="./values.yaml#L52">global.databaseUsername</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django_db_user"
</pre>
</div>
			</td>
			<td>Database username backend to connect to. If you use external backend, provide the value</td>
		</tr>
		<tr>
			<td id="global--debug"><a href="./values.yaml#L67">global.debug</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"False"
</pre>
</div>
			</td>
			<td>Python boolean literal, this will correspond to `DEBUG` environment variable inside the Django container. Useful as a debug switch.</td>
		</tr>
		<tr>
			<td id="global--djangoArgs"><a href="./values.yaml#L30">global.djangoArgs</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
global.djangoArgs: |
  ["uwsgi","--chdir=${REPO_ROOT}","--module=${MAIN_APP_NAME}.wsgi","--socket=:8000","--http=0.0.0.0:8080","--processes=5","--buffer-size=8192"]
 
</pre>
</div>
			</td>
			<td>The django command args to be passed to entrypoint command</td>
		</tr>
		<tr>
			<td id="global--djangoCommand"><a href="./values.yaml#L26">global.djangoCommand</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
global.djangoCommand: |
  ["/opt/django/scripts/docker-entrypoint.sh"]
 
</pre>
</div>
			</td>
			<td>The django entrypoint command to execute</td>
		</tr>
		<tr>
			<td id="global--djangoSecretKey--value"><a href="./values.yaml#L46">global.djangoSecretKey.value</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
null
</pre>
</div>
			</td>
			<td>Specify this Django Secret string value. If not, it will be autogenerated everytime chart upgraded</td>
		</tr>
		<tr>
			<td id="global--djangoSettingsModule"><a href="./values.yaml#L71">global.djangoSettingsModule</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django.settings"
</pre>
</div>
			</td>
			<td>Django settings module to be used</td>
		</tr>
		<tr>
			<td id="global--existingSecret"><a href="./values.yaml#L18">global.existingSecret</a></td>
			<td>
tpl/string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
global.existingSecret: |
 
</pre>
</div>
			</td>
			<td>Name of existing secret</td>
		</tr>
		<tr>
			<td id="global--mainAppName"><a href="./values.yaml#L69">global.mainAppName</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django"
</pre>
</div>
			</td>
			<td>The main app name to execute. Affects which settings, wsgi, and rootURL to use.</td>
		</tr>
		<tr>
			<td id="global--mediaRoot"><a href="./values.yaml#L77">global.mediaRoot</a></td>
			<td>
path
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"/opt/django/media"
</pre>
</div>
			</td>
			<td>Location to the media directory</td>
		</tr>
		<tr>
			<td id="global--rootURLConf"><a href="./values.yaml#L73">global.rootURLConf</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django.urls"
</pre>
</div>
			</td>
			<td>Django root URL conf to be used</td>
		</tr>
		<tr>
			<td id="global--sharedSecretName"><a href="./values.yaml#L20">global.sharedSecretName</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django-shared-secret"
</pre>
</div>
			</td>
			<td>Name of shared secret store that will be generated</td>
		</tr>
		<tr>
			<td id="global--siteName"><a href="./values.yaml#L23">global.siteName</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django"
</pre>
</div>
			</td>
			<td>The site name. It will be used to construct url such as http://django/</td>
		</tr>
		<tr>
			<td id="global--staticRoot"><a href="./values.yaml#L75">global.staticRoot</a></td>
			<td>
path
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"/opt/django/static"
</pre>
</div>
			</td>
			<td>Location to the static directory</td>
		</tr>
		<tr>
			<td id="image"><a href="./values.yaml#L2">image</a></td>
			<td>
object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{
  "pullPolicy": "IfNotPresent",
  "registry": "docker.io",
  "repository": "lucernae/django-sample",
  "tag": "3.1"
}
</pre>
</div>
			</td>
			<td>Image map</td>
		</tr>
		<tr>
			<td id="image--pullPolicy"><a href="./values.yaml#L10">image.pullPolicy</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"IfNotPresent"
</pre>
</div>
			</td>
			<td>Image pullPolicy</td>
		</tr>
		<tr>
			<td id="image--registry"><a href="./values.yaml#L4">image.registry</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"docker.io"
</pre>
</div>
			</td>
			<td>Image registry</td>
		</tr>
		<tr>
			<td id="image--repository"><a href="./values.yaml#L6">image.repository</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"lucernae/django-sample"
</pre>
</div>
			</td>
			<td>Image repository</td>
		</tr>
		<tr>
			<td id="image--tag"><a href="./values.yaml#L8">image.tag</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"3.1"
</pre>
</div>
			</td>
			<td>Image tag</td>
		</tr>
		<tr>
			<td id="ingress--annotations"><a href="./values.yaml#L155">ingress.annotations</a></td>
			<td>
dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td>Custom Ingress annotations</td>
		</tr>
		<tr>
			<td id="ingress--enabled"><a href="./values.yaml#L150">ingress.enabled</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
false
</pre>
</div>
			</td>
			<td>Set to true to generate Ingress resource</td>
		</tr>
		<tr>
			<td id="ingress--host"><a href="./values.yaml#L153">ingress.host</a></td>
			<td>
tpl/string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
ingress.host: |
 
</pre>
</div>
			</td>
			<td>Set custom host name. (DNS name convention)</td>
		</tr>
		<tr>
			<td id="ingress--labels"><a href="./values.yaml#L157">ingress.labels</a></td>
			<td>
dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td>Custom Ingress labels</td>
		</tr>
		<tr>
			<td id="ingress--tls--enabled"><a href="./values.yaml#L160">ingress.tls.enabled</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
false
</pre>
</div>
			</td>
			<td>Set to true to enable HTTPS</td>
		</tr>
		<tr>
			<td id="ingress--tls--secretName"><a href="./values.yaml#L162">ingress.tls.secretName</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"django-tls"
</pre>
</div>
			</td>
			<td>You must provide a secret name where the TLS cert is stored</td>
		</tr>
		<tr>
			<td id="labels"><a href="./values.yaml#L81">labels</a></td>
			<td>
map
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="yaml">
user/workload: "true"
client-name: "my-boss"
project-name: "awesome-project"

</pre>
</div>
			</td>
			<td>The deployment label</td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--accessModes[0]"><a href="./values.yaml#L199">persistence.mediaDir.accessModes[0]</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"ReadWriteOnce"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--annotations"><a href="./values.yaml#L200">persistence.mediaDir.annotations</a></td>
			<td>
object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--enabled"><a href="./values.yaml#L193">persistence.mediaDir.enabled</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
true
</pre>
</div>
			</td>
			<td>Allow persistence</td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--existingClaim"><a href="./values.yaml#L194">persistence.mediaDir.existingClaim</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
false
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--mountPath"><a href="./values.yaml#L195">persistence.mediaDir.mountPath</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"/opt/django/media"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--size"><a href="./values.yaml#L197">persistence.mediaDir.size</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"8Gi"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--subPath"><a href="./values.yaml#L196">persistence.mediaDir.subPath</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"media"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--accessModes"><a href="./values.yaml#L188">persistence.staticDir.accessModes</a></td>
			<td>
<a target="_blank"
   href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes"
   >k8s/storage/persistent-volume/access-modes</a>
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="yaml">
- ReadWriteOnce

</pre>
</div>
			</td>
			<td>Static Dir access modes</td>
		</tr>
		<tr>
			<td id="persistence--staticDir--annotations"><a href="./values.yaml#L190">persistence.staticDir.annotations</a></td>
			<td>
object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--enabled"><a href="./values.yaml#L181">persistence.staticDir.enabled</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
true
</pre>
</div>
			</td>
			<td>Allow persistence</td>
		</tr>
		<tr>
			<td id="persistence--staticDir--existingClaim"><a href="./values.yaml#L182">persistence.staticDir.existingClaim</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
false
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--mountPath"><a href="./values.yaml#L183">persistence.staticDir.mountPath</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"/opt/django/static"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--size"><a href="./values.yaml#L185">persistence.staticDir.size</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"8Gi"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--subPath"><a href="./values.yaml#L184">persistence.staticDir.subPath</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"static"
</pre>
</div>
			</td>
			<td></td>
		</tr>
		<tr>
			<td id="postgis--enabled"><a href="./values.yaml#L170">postgis.enabled</a></td>
			<td>
bool
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
true
</pre>
</div>
			</td>
			<td>Enable postgis as database backend by default. Set to false if using different external backend.</td>
		</tr>
		<tr>
			<td id="postgis--existingSecret"><a href="./values.yaml#L174">postgis.existingSecret</a></td>
			<td>
tpl/string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
postgis.existingSecret: |
  {{ include "common.sharedSecretName" . | quote -}}
 
</pre>
</div>
			</td>
			<td>Existing secret to be used</td>
		</tr>
		<tr>
			<td id="probe"><a href="./values.yaml#L166">probe</a></td>
			<td>
tpl/object
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
probe: |
 
</pre>
</div>
			</td>
			<td>Probe can be overridden</td>
		</tr>
		<tr>
			<td id="sampleYaml"><a href="./values.yaml#L227">sampleYaml</a></td>
			<td>
dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td>Values with long description
Sometimes you need a very long description
for your values.

Any comment section for a given key with **@raw** attribute
will be treated as raw string and stored as is.
Since it generates in Markdown format, you can do something like this:

```yaml
hello:
  bar: true
```

Markdown also accept subset of HTML tags. So you can also do this:

<details>
<summary>+Expand</summary>

```bash
execute some command
```

</details></td>
		</tr>
		<tr>
			<td id="service--annotations"><a href="./values.yaml#L146">service.annotations</a></td>
			<td>
dict
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
{}
</pre>
</div>
			</td>
			<td>Extra service annotations</td>
		</tr>
		<tr>
			<td id="service--clusterIP"><a href="./values.yaml#L135">service.clusterIP</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
""
</pre>
</div>
			</td>
			<td>Specify `None` for headless service. Otherwise, leave them be.</td>
		</tr>
		<tr>
			<td id="service--externalIPs"><a href="./values.yaml#L138">service.externalIPs</a></td>
			<td>
tpl/array
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="tpl">
service.externalIPs: |
 
</pre>
</div>
			</td>
			<td>Specify for LoadBalancer service type</td>
		</tr>
		<tr>
			<td id="service--nodePort"><a href="./values.yaml#L143">service.nodePort</a></td>
			<td>
int
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
null
</pre>
</div>
			</td>
			<td>Specify node port, for NodePort service type</td>
		</tr>
		<tr>
			<td id="service--port"><a href="./values.yaml#L140">service.port</a></td>
			<td>
int
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
80
</pre>
</div>
			</td>
			<td>Specify service port</td>
		</tr>
		<tr>
			<td id="service--type"><a href="./values.yaml#L133">service.type</a></td>
			<td>
string
</td>
			<td>
				<div style="max-width: 300px;">
<pre lang="json">
"ClusterIP"
</pre>
</div>
			</td>
			<td>Define k8s service for Django.</td>
		</tr>
	</tbody>
</table>

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.10.0](https://github.com/norwoodj/helm-docs/releases/v1.10.0)
