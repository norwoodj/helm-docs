package helm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseComment_WithExample(t *testing.T) {
	commentLines := []string{
		"# foo.bar -- This is bar",
		"# @example -- `{\"some\": \"thing\", \"very\": \"complex\"}`",
	}
	key, c := ParseComment(commentLines)
	assert.Equal(t, "foo.bar", key)
	assert.Equal(t, "This is bar", c.Description)
	assert.Equal(t, "`{\"some\": \"thing\", \"very\": \"complex\"}`", c.Example)
}

func TestParseComment_WithExampleAndDefault(t *testing.T) {
	commentLines := []string{
		"# baz -- This is baz",
		"# @default -- \"\"",
		"# @example -- Accepted values: `john`, `doe`",
	}
	key, c := ParseComment(commentLines)
	assert.Equal(t, "baz", key)
	assert.Equal(t, "This is baz", c.Description)
	assert.Equal(t, "\"\"", c.Default)
	assert.Equal(t, "Accepted values: `john`, `doe`", c.Example)
}

func TestParseComment_WithoutExample(t *testing.T) {
	commentLines := []string{
		"# quz -- No example for quz",
	}
	key, c := ParseComment(commentLines)
	assert.Equal(t, "quz", key)
	assert.Equal(t, "No example for quz", c.Description)
	assert.Equal(t, "", c.Example)
}
