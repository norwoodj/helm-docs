# django

![Version: 0.2.1](https://img.shields.io/badge/Version-0.2.1-informational?style=flat-square) ![AppVersion: 3.1](https://img.shields.io/badge/AppVersion-3.1-informational?style=flat-square)

Generic chart for basic Django-based web app

**Homepage:** <https://www.djangoproject.com/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Rizky Maulana Nugraha | lana.pcfre@gmail.com |  |

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
want to combine helm-docs with a JAMStack static generator where you can have your own page generated
from this README.

The customization can goes even further. Normally, you can't define anchor in markdown unless it it a heading. But you can do so easily using HTML tags.
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
			<td id="extraConfigMap">extraConfigMap</td>
			<td>
tpl/dict
</td>
			<td><pre lang="tpl">
extraConfigMap: |
 
</pre>
</td>
			<td>Define this for extra config map to be included in django-shared-config</td>
		</tr>
		<tr>
			<td id="extraPodEnv">extraPodEnv</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
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
</td>
			<td>Define this for extra Django environment variables</td>
		</tr>
		<tr>
			<td id="extraPodSpec">extraPodSpec</td>
			<td>
tpl/object
</td>
			<td><pre lang="tpl">
extraPodSpec: |
 
</pre>
</td>
			<td>This will be evaluated as pod spec</td>
		</tr>
		<tr>
			<td id="extraSecret">extraSecret</td>
			<td>
tpl/dict
</td>
			<td><pre lang="tpl">
extraSecret: |
 
</pre>
</td>
			<td>Define this for extra secrets to be included in django-shared-secret secret</td>
		</tr>
		<tr>
			<td id="extraVolume">extraVolume</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
extraVolume: |
 
</pre>
</td>
			<td>Define this for extra volume (in pair with extraVolumeMounts)</td>
		</tr>
		<tr>
			<td id="extraVolumeMounts">extraVolumeMounts</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
extraVolumeMounts: |
 
</pre>
</td>
			<td>Define this for extra volume mounts in the pod</td>
		</tr>
		<tr>
			<td id="global">global</td>
			<td>
object
</td>
			<td><pre lang="json">
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
</td>
			<td>This key name is used for service interconnection between subcharts and parent charts.</td>
		</tr>
		<tr>
			<td id="global--adminEmail">global.adminEmail</td>
			<td>
<a href="#stringemail" title="
This value type is for a valid email address format. Such as owner@somedomain.org.">string/email</a>
</td>
			<td><pre lang="email">
<a href="mailto:admin@localhost">"admin@localhost"</a>
</pre>
</td>
			<td>Default admin email sender</td>
		</tr>
		<tr>
			<td id="global--adminPassword--value">global.adminPassword.value</td>
			<td>
string
</td>
			<td><pre lang="json">
null
</pre>
</td>
			<td>Specify this password value. If not, it will be autogenerated everytime chart upgraded</td>
		</tr>
		<tr>
			<td id="global--adminUser">global.adminUser</td>
			<td>
string
</td>
			<td><pre lang="json">
"admin"
</pre>
</td>
			<td>Default super user admin username</td>
		</tr>
		<tr>
			<td id="global--databaseHost">global.databaseHost</td>
			<td>
string
</td>
			<td><pre lang="json">
"postgis"
</pre>
</td>
			<td>Django database host location. By default this chart can generate standard postgres chart. So you can leave it as default. If you use external backend,  you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databaseName">global.databaseName</td>
			<td>
string
</td>
			<td><pre lang="json">
"django"
</pre>
</td>
			<td>Django database name</td>
		</tr>
		<tr>
			<td id="global--databasePassword--value">global.databasePassword.value</td>
			<td>
string
</td>
			<td><pre lang="json">
null
</pre>
</td>
			<td>Specify this password value. If not, it will be autogenerated everytime chart upgraded. If you use external backend, you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databasePort">global.databasePort</td>
			<td>
int
</td>
			<td><pre lang="json">
5432
</pre>
</td>
			<td>Django database port. By default this chart can generate standard postgres chart. So you can leave it as default. If you use external backend,  you must provide the value</td>
		</tr>
		<tr>
			<td id="global--databaseUsername">global.databaseUsername</td>
			<td>
string
</td>
			<td><pre lang="json">
"django_db_user"
</pre>
</td>
			<td>Database username backend to connect to. If you use external backend, provide the value</td>
		</tr>
		<tr>
			<td id="global--debug">global.debug</td>
			<td>
string
</td>
			<td><pre lang="json">
"False"
</pre>
</td>
			<td>Python boolean literal, this will correspond to `DEBUG` environment variable inside the Django container. Useful as a debug switch.</td>
		</tr>
		<tr>
			<td id="global--djangoArgs">global.djangoArgs</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
global.djangoArgs: |
  ["uwsgi","--chdir=${REPO_ROOT}","--module=${MAIN_APP_NAME}.wsgi","--socket=:8000","--http=0.0.0.0:8080","--processes=5","--buffer-size=8192"]
 
</pre>
</td>
			<td>The django command args to be passed to entrypoint command</td>
		</tr>
		<tr>
			<td id="global--djangoCommand">global.djangoCommand</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
global.djangoCommand: |
  ["/opt/django/scripts/docker-entrypoint.sh"]
 
</pre>
</td>
			<td>The django entrypoint command to execute</td>
		</tr>
		<tr>
			<td id="global--djangoSecretKey--value">global.djangoSecretKey.value</td>
			<td>
string
</td>
			<td><pre lang="json">
null
</pre>
</td>
			<td>Specify this Django Secret string value. If not, it will be autogenerated everytime chart upgraded</td>
		</tr>
		<tr>
			<td id="global--djangoSettingsModule">global.djangoSettingsModule</td>
			<td>
string
</td>
			<td><pre lang="json">
"django.settings"
</pre>
</td>
			<td>Django settings module to be used</td>
		</tr>
		<tr>
			<td id="global--existingSecret">global.existingSecret</td>
			<td>
tpl/string
</td>
			<td><pre lang="tpl">
global.existingSecret: |
 
</pre>
</td>
			<td>Name of existing secret</td>
		</tr>
		<tr>
			<td id="global--mainAppName">global.mainAppName</td>
			<td>
string
</td>
			<td><pre lang="json">
"django"
</pre>
</td>
			<td>The main app name to execute. Affects which settings, wsgi, and rootURL to use.</td>
		</tr>
		<tr>
			<td id="global--mediaRoot">global.mediaRoot</td>
			<td>
path
</td>
			<td><pre lang="json">
"/opt/django/media"
</pre>
</td>
			<td>Location to the media directory</td>
		</tr>
		<tr>
			<td id="global--rootURLConf">global.rootURLConf</td>
			<td>
string
</td>
			<td><pre lang="json">
"django.urls"
</pre>
</td>
			<td>Django root URL conf to be used</td>
		</tr>
		<tr>
			<td id="global--sharedSecretName">global.sharedSecretName</td>
			<td>
string
</td>
			<td><pre lang="json">
"django-shared-secret"
</pre>
</td>
			<td>Name of shared secret store that will be generated</td>
		</tr>
		<tr>
			<td id="global--staticRoot">global.staticRoot</td>
			<td>
path
</td>
			<td><pre lang="json">
"/opt/django/static"
</pre>
</td>
			<td>Location to the static directory</td>
		</tr>
		<tr>
			<td id="image">image</td>
			<td>
object
</td>
			<td><pre lang="json">
{
  "pullPolicy": "IfNotPresent",
  "registry": "docker.io",
  "repository": "lucernae/django-sample",
  "tag": "3.1"
}
</pre>
</td>
			<td>Image map</td>
		</tr>
		<tr>
			<td id="image--pullPolicy">image.pullPolicy</td>
			<td>
string
</td>
			<td><pre lang="json">
"IfNotPresent"
</pre>
</td>
			<td>Image pullPolicy</td>
		</tr>
		<tr>
			<td id="image--registry">image.registry</td>
			<td>
string
</td>
			<td><pre lang="json">
"docker.io"
</pre>
</td>
			<td>Image registry</td>
		</tr>
		<tr>
			<td id="image--repository">image.repository</td>
			<td>
string
</td>
			<td><pre lang="json">
"lucernae/django-sample"
</pre>
</td>
			<td>Image repository</td>
		</tr>
		<tr>
			<td id="image--tag">image.tag</td>
			<td>
string
</td>
			<td><pre lang="json">
"3.1"
</pre>
</td>
			<td>Image tag</td>
		</tr>
		<tr>
			<td id="ingress--annotations">ingress.annotations</td>
			<td>
dict
</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Custom Ingress annotations</td>
		</tr>
		<tr>
			<td id="ingress--enabled">ingress.enabled</td>
			<td>
bool
</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>Set to true to generate Ingress resource</td>
		</tr>
		<tr>
			<td id="ingress--host">ingress.host</td>
			<td>
tpl/string
</td>
			<td><pre lang="tpl">
ingress.host: |
 
</pre>
</td>
			<td>Set custom host name. (DNS name convention)</td>
		</tr>
		<tr>
			<td id="ingress--labels">ingress.labels</td>
			<td>
dict
</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Custom Ingress labels</td>
		</tr>
		<tr>
			<td id="ingress--tls--enabled">ingress.tls.enabled</td>
			<td>
bool
</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td>Set to true to enable HTTPS</td>
		</tr>
		<tr>
			<td id="ingress--tls--secretName">ingress.tls.secretName</td>
			<td>
string
</td>
			<td><pre lang="json">
"django-tls"
</pre>
</td>
			<td>You must provide a secret name where the TLS cert is stored</td>
		</tr>
		<tr>
			<td id="labels">labels</td>
			<td>
map
</td>
			<td><pre lang="yaml">
user/workload: "true"
client-name: "my-boss"
project-name: "awesome-project"

</pre>
</td>
			<td>The deployment label</td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--accessModes[0]">persistence.mediaDir.accessModes[0]</td>
			<td>
string
</td>
			<td><pre lang="json">
"ReadWriteOnce"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--annotations">persistence.mediaDir.annotations</td>
			<td>
object
</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--enabled">persistence.mediaDir.enabled</td>
			<td>
bool
</td>
			<td><pre lang="json">
true
</pre>
</td>
			<td>Allow persistence</td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--existingClaim">persistence.mediaDir.existingClaim</td>
			<td>
bool
</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--mountPath">persistence.mediaDir.mountPath</td>
			<td>
string
</td>
			<td><pre lang="json">
"/opt/django/media"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--size">persistence.mediaDir.size</td>
			<td>
string
</td>
			<td><pre lang="json">
"8Gi"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--mediaDir--subPath">persistence.mediaDir.subPath</td>
			<td>
string
</td>
			<td><pre lang="json">
"media"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--accessModes">persistence.staticDir.accessModes</td>
			<td>
<a target="_blank"
   href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes"
   >k8s/storage/persistent-volume/access-modes</a>
</td>
			<td><pre lang="yaml">
- ReadWriteOnce

</pre>
</td>
			<td>Static Dir access modes</td>
		</tr>
		<tr>
			<td id="persistence--staticDir--annotations">persistence.staticDir.annotations</td>
			<td>
object
</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--enabled">persistence.staticDir.enabled</td>
			<td>
bool
</td>
			<td><pre lang="json">
true
</pre>
</td>
			<td>Allow persistence</td>
		</tr>
		<tr>
			<td id="persistence--staticDir--existingClaim">persistence.staticDir.existingClaim</td>
			<td>
bool
</td>
			<td><pre lang="json">
false
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--mountPath">persistence.staticDir.mountPath</td>
			<td>
string
</td>
			<td><pre lang="json">
"/opt/django/static"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--size">persistence.staticDir.size</td>
			<td>
string
</td>
			<td><pre lang="json">
"8Gi"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="persistence--staticDir--subPath">persistence.staticDir.subPath</td>
			<td>
string
</td>
			<td><pre lang="json">
"static"
</pre>
</td>
			<td></td>
		</tr>
		<tr>
			<td id="postgis--enabled">postgis.enabled</td>
			<td>
bool
</td>
			<td><pre lang="json">
true
</pre>
</td>
			<td>Enable postgis as database backend by default. Set to false if using different external backend.</td>
		</tr>
		<tr>
			<td id="postgis--existingSecret">postgis.existingSecret</td>
			<td>
tpl/string
</td>
			<td><pre lang="tpl">
postgis.existingSecret: |
  {{ include "common.sharedSecretName" . | quote -}}
 
</pre>
</td>
			<td>Existing secret to be used</td>
		</tr>
		<tr>
			<td id="probe">probe</td>
			<td>
tpl/object
</td>
			<td><pre lang="tpl">
probe: |
 
</pre>
</td>
			<td>Probe can be overridden</td>
		</tr>
		<tr>
			<td id="service--annotations">service.annotations</td>
			<td>
dict
</td>
			<td><pre lang="json">
{}
</pre>
</td>
			<td>Extra service annotations</td>
		</tr>
		<tr>
			<td id="service--clusterIP">service.clusterIP</td>
			<td>
string
</td>
			<td><pre lang="json">
""
</pre>
</td>
			<td>Specify `None` for headless service. Otherwise, leave them be.</td>
		</tr>
		<tr>
			<td id="service--externalIPs">service.externalIPs</td>
			<td>
tpl/array
</td>
			<td><pre lang="tpl">
service.externalIPs: |
 
</pre>
</td>
			<td>Specify for LoadBalancer service type</td>
		</tr>
		<tr>
			<td id="service--nodePort">service.nodePort</td>
			<td>
int
</td>
			<td><pre lang="json">
null
</pre>
</td>
			<td>Specify node port, for NodePort service type</td>
		</tr>
		<tr>
			<td id="service--port">service.port</td>
			<td>
int
</td>
			<td><pre lang="json">
80
</pre>
</td>
			<td>Specify service port</td>
		</tr>
		<tr>
			<td id="service--type">service.type</td>
			<td>
string
</td>
			<td><pre lang="json">
"ClusterIP"
</pre>
</td>
			<td>Define k8s service for Django.</td>
		</tr>
	</tbody>
</table>

