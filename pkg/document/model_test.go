package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestGetUnsortedValueRowsEmptyValues(t *testing.T) {
	tests := []struct {
		name   string
		values string
	}{
		{"empty file", ""},
		{"comments only", "# No configurable values\n"},
		{"document marker", "---\n"},
		{"document with comments", "---\n# No configurable values\n"},
		{"document with end marker", "---\n...\n"},
		{"explicit null", "null\n"},
		{"null shorthand", "~\n"},
		{"empty map", "{}\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var values yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(tt.values), &values))

			rows, err := getUnsortedValueRows(&values, nil)
			require.NoError(t, err)
			assert.Empty(t, rows)
		})
	}
}

func TestGetUnsortedValueRowsRejectsNonMappingValues(t *testing.T) {
	for _, values := range []string{"[]", "[one, two]", "text", `""`, `"null"`, "false", "0"} {
		t.Run(values, func(t *testing.T) {
			var document yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(values), &document))

			rows, err := getUnsortedValueRows(&document, nil)
			require.ErrorContains(t, err, "values file must resolve to a map")
			assert.Empty(t, rows)
		})
	}
}
