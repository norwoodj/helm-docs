package helm_test

import (
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/stretchr/testify/assert"
)

func TestParseComment(t *testing.T) {
	commentLines := []string{
		"# controller.image.repository -- The repository of the controller image",
	}
	valueKey, c := helm.ParseComment(commentLines)

	assert.Equal(t, "controller.image.repository", valueKey)
	assert.Equal(t, "The repository of the controller image", c.Description)
}

func TestParseCommentNewStyle(t *testing.T) {
	commentLines := []string{
		"# -- The repository of the controller image",
	}
	valueKey, c := helm.ParseComment(commentLines)

	assert.Equal(t, "", valueKey)
	assert.Equal(t, "The repository of the controller image", c.Description)
}
