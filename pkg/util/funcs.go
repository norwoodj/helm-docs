package util

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

func FuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	f["toYaml"] = toYAML
	f["fromYaml"] = fromYAML
	f["toMdTableCell"] = toMdTableCell
	f["toMdDefault"] = toMdDefault
	f["toHtmlTableCell"] = toHtmlTableCell
	f["yamlDefault"] = yamlDefault
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

var codeFenceRegex = regexp.MustCompile("(?s)```(\\S*)\n(.*?)\n?```")

// toMdTableCell makes a multiline description fit in one markdown table cell:
// ``` fences become <pre> blocks and the remaining newlines become <br>.
func toMdTableCell(s string) string { return toTableCell(s, true) }

// toHtmlTableCell is toMdTableCell for an HTML table, where markdown escapes like \| are not processed.
func toHtmlTableCell(s string) string { return toTableCell(s, false) }

func toTableCell(s string, escapePipe bool) string {
	s = codeFenceRegex.ReplaceAllStringFunc(s, func(block string) string {
		m := codeFenceRegex.FindStringSubmatch(block)
		code := html.EscapeString(m[2])
		if escapePipe {
			// a bare | (e.g. YAML block scalar) would end the markdown table cell
			code = strings.ReplaceAll(code, "|", "\\|")
		}
		if m[1] == "" {
			return "<pre>" + code + "</pre>"
		}
		return fmt.Sprintf(`<pre lang="%s">%s</pre>`, html.EscapeString(m[1]), code)
	})
	return strings.ReplaceAll(s, "\n", "<br>")
}

// toMdDefault renders a non-empty object/list default (`<json>`) as a YAML <pre> block for a markdown table cell.
// Anything else (scalars, {}, [], non-backticked text) is returned unchanged.
func toMdDefault(s string) string {
	if y := yamlDefault(s); y != "" {
		return toMdTableCell("```yaml\n" + y + "\n```")
	}
	return s
}

// yamlDefault returns a non-empty object/list default (`<json>`) as YAML, or "" for anything else.
func yamlDefault(s string) string {
	if len(s) < 2 || !strings.HasPrefix(s, "`") || !strings.HasSuffix(s, "`") {
		return ""
	}
	// JSON is valid YAML; decoding it as YAML keeps big ints as ints (no 1e+06)
	var v interface{}
	if err := yaml.Unmarshal([]byte(s[1:len(s)-1]), &v); err != nil {
		return ""
	}
	switch t := v.(type) {
	case map[string]interface{}:
		if len(t) == 0 {
			return ""
		}
	case []interface{}:
		if len(t) == 0 {
			return ""
		}
	default:
		return ""
	}
	var b strings.Builder
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return ""
	}
	return strings.TrimSuffix(b.String(), "\n")
}
