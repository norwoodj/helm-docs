package document_test

import (
	"path/filepath"
	"testing"

	"github.com/norwoodj/helm-docs/pkg/helm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	. "github.com/norwoodj/helm-docs/pkg/document"
)

func TestGetDependencyValues(t *testing.T) {
	type args struct {
		root helm.ChartDocumentationInfo
		all  []helm.ChartDocumentationInfo
	}
	tests := []struct {
		name    string
		args    args
		want    []DependencyValues
		wantErr bool
	}{
		{
			name: "zero dependencies",
		},
		{
			name: "local dependency with name",
			args: args{
				root: info([]helm.ChartDependenciesItem{{Name: "sub-name"}}, "root"),
				all: []helm.ChartDocumentationInfo{
					info(nil, "root", "charts", "sub-name"),
				},
			},
			want: []DependencyValues{
				values("sub-name", "root", "charts", "sub-name"),
			},
		},
		{
			name: "local dependency with alias",
			args: args{
				root: info([]helm.ChartDependenciesItem{{Name: "sub-name", Alias: "sub-alias"}}, "root"),
				all: []helm.ChartDocumentationInfo{
					info(nil, "root", "charts", "sub-name"),
				},
			},
			want: []DependencyValues{
				values("sub-alias", "root", "charts", "sub-name"),
			},
		},
		{
			name: "nested dependencies",
			args: args{
				root: info([]helm.ChartDependenciesItem{{Name: "mid-a"}, {Name: "mid-b"}}, "root"),
				all: []helm.ChartDocumentationInfo{
					info([]helm.ChartDependenciesItem{{Name: "leaf-c"}, {Name: "leaf-d"}}, "root", "charts", "mid-a"),
					info([]helm.ChartDependenciesItem{{Name: "leaf-e"}, {Name: "leaf-f"}}, "root", "charts", "mid-b"),
					info(nil, "root", "charts", "mid-a", "charts", "leaf-c"),
					info(nil, "root", "charts", "mid-a", "charts", "leaf-d"),
					info(nil, "root", "charts", "mid-b", "charts", "leaf-e"),
					info(nil, "root", "charts", "mid-b", "charts", "leaf-f"),
				},
			},
			want: []DependencyValues{
				values("mid-a", "root", "charts", "mid-a"),
				values("mid-b", "root", "charts", "mid-b"),
				values("mid-a.leaf-c", "root", "charts", "mid-a", "charts", "leaf-c"),
				values("mid-a.leaf-d", "root", "charts", "mid-a", "charts", "leaf-d"),
				values("mid-b.leaf-e", "root", "charts", "mid-b", "charts", "leaf-e"),
				values("mid-b.leaf-f", "root", "charts", "mid-b", "charts", "leaf-f"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			infosByChartPath := make(map[string]helm.ChartDocumentationInfo)
			for _, info := range tt.args.all {
				infosByChartPath[info.ChartDirectory] = info
			}
			got, err := GetDependencyValues(tt.args.root, infosByChartPath)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func info(dependencies []helm.ChartDependenciesItem, dirParts ...string) helm.ChartDocumentationInfo {
	dir := filepath.Join(dirParts...)
	return helm.ChartDocumentationInfo{
		ChartDirectory:          dir,
		ChartValues:             &yaml.Node{Value: dir},
		ChartValuesDescriptions: map[string]helm.ChartValueDescription{"value": {Description: dir}},
		ChartDependencies: helm.ChartDependencies{
			Dependencies: dependencies,
		},
	}
}

func values(prefix string, dirParts ...string) DependencyValues {
	dir := filepath.Join(dirParts...)
	return DependencyValues{
		Prefix:                  prefix,
		ChartValues:             &yaml.Node{Value: dir},
		ChartValuesDescriptions: map[string]helm.ChartValueDescription{"value": {Description: dir}},
	}
}
