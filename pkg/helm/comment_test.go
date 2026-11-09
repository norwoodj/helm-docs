package helm_test

import (
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/stretchr/testify/assert"
)

func TestParseCommentStopsAtExtra(t *testing.T) {
	key, c := helm.ParseComment([]string{
		"# -- List of persistent volume claims to create.",
		"# For very long comments, break them into multiple lines.",
		"# @default -- the chart will construct this list internally unless specified",
		"# @extra -- Volume Claim Template",
		"# The volume claim template is used to create a `PersistentVolumeClaim`.",
	})
	assert.Equal(t, "", key)
	assert.Equal(t, "List of persistent volume claims to create. For very long comments, break them into multiple lines.", c.Description)
	assert.Equal(t, "the chart will construct this list internally unless specified", c.Default)
}

func TestParseCommentEmpty(t *testing.T) {
	key, c := helm.ParseComment(nil)
	assert.Equal(t, "", key)
	assert.Equal(t, helm.ChartValueDescription{}, c)
}
