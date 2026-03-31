package util

import (
	"html"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	f["toYaml"] = toYAML
	f["fromYaml"] = fromYAML
	f["htmlEscape"] = htmlEscape
	return f
}

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// fromYAML converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// htmlEscape escapes special HTML characters in a string to their HTML entity equivalents.
// It also converts Unicode escape sequences (\u003c, \u003e, \u0026) produced by Go's json.Marshal
// to their HTML entity equivalents (&lt;, &gt;, &amp;).
//
// This is necessary because Sprig's toPrettyJson function uses json.MarshalIndent without
// SetEscapeHTML(false), which means it escapes <, >, and & to Unicode sequences.
// We want proper HTML entities instead for better readability in markdown/HTML output.
//
// This is designed to be called from a template.
func htmlEscape(s string) string {
	// First, replace Unicode escape sequences with actual characters
	s = strings.ReplaceAll(s, `\u003c`, "<")
	s = strings.ReplaceAll(s, `\u003e`, ">")
	s = strings.ReplaceAll(s, `\u0026`, "&")

	// Then apply HTML escaping to convert them to HTML entities
	return html.EscapeString(s)
}
