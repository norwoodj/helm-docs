# Documenting empty maps and lists

To document values that are not present in the `values.yaml` file, you can use the `@parseFoot` annotation.
This will parse the foot comment as a yaml, treating it as it was part of the `values.yaml` file, allowing you to document values that are empty by default.

Requirements and limitations:
- The `@parseFoot` annotation must be on the object or list.
- The foot comment must be valid yaml. Indentation before or after `#` does not matter and leading whitespace will be removed, it just needs to be on a same level.
- There needs to be a blank line between the end of the foot comment and head comment of next value.
- The keys for commented values must not be present, or the description will not work.
  ex: `env[].name -- Env name` is not valid, needs to be `-- Env name`
- You can use one item to document a list, the `[0]` will be replaced with `[]`, but `[1]` will be kept as is.
- You can use any key to document a map, it's kept as is, so you can even have multiple items to document different use cases, but avoid using `"*"` without escaping it, as it can be rendered as italic in markdown in double nested item, prefer using ``"`*`"`` instead.

## Values

### List

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| envs | list | `[]` | Environment variables |
| envs[1] | object | `{"name":"weird item"}` | Avoid this |
| envs[1].name | string | `"weird item"` | A bad item, should be all in one item |
| envs[].name | string | `""` | Environment variable name |
| envs[].value | string | `""` | Environment variable value |
| envs[].valueFrom | object | `nil` | Environment variable supplied from secret or configmap |
| envs[].valueFrom.secretKeyRef | object | `nil` | Environment variable supplied from secret |
| envs[].valueFrom.secretKeyRef.key | string | `""` | Key to mount to environment variable |
| envs[].valueFrom.secretKeyRef.name | string | `""` | Kubernetes Secret Name |

### Map

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| map | object | `{}` | Map with dynamic keys, with object children |
| map.* | object | `{"incorrectNestedMap":{},"nestedKey":"value","nestedMap":{}}` | A map item |
| map.*.incorrectNestedMap | object | `{}` | A nested dynamic map, it can have any keys. This is a map of int |
| map.*.incorrectNestedMap.* | int | `nil` | A dynamic map item - this will cause italic in rendering because neither asterisk is in code tags |
| map.*.nestedKey | string | `"value"` | A nested map item |
| map.*.nestedMap | object | `{}` | A nested dynamic map, showing multiple levels of foot comment parsing. This is a map of string |
| map.*.nestedMap.`*` | string | `"value"` | A dynamic map item |
| map.[specific-example] | object | `{"differentExlusiveKey":"value"}` | A specific example key, it can be rendered differently in documentation if specified like this |
| map.[specific-example].differentExlusiveKey | string | `"value"` | A specific example key item |
