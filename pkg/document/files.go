package document

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"gopkg.in/yaml.v3"
)

// Near identical to https://github.com/helm/helm/blob/main/pkg/engine/files.go as to preserve the interface.

type files map[string][]byte

func getFiles(dir string) (files, error) {
	result := make(files)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		result[path] = data

		return nil
	})
	if err != nil {
		return map[string][]byte{}, err
	}

	return result, nil
}

func (f files) GetBytes(name string) []byte {
	if v, ok := f[name]; ok {
		return v
	}
	return []byte{}
}

func (f files) Get(name string) string {
	return string(f.GetBytes(name))
}

func (f files) Glob(pattern string) files {
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		g, _ = glob.Compile("**")
	}

	result := make(files)
	for name, contents := range f {
		if g.Match(name) {
			result[name] = contents
		}
	}

	return result
}

func (f files) AsConfig() string {
	if f == nil {
		return ""
	}

	m := make(map[string]string)

	// Explicitly convert to strings, and file names
	for k, v := range f {
		m[path.Base(k)] = string(v)
	}

	return toYAML(m)
}

func (f files) AsSecrets() string {
	if f == nil {
		return ""
	}

	m := make(map[string]string)

	for k, v := range f {
		m[path.Base(k)] = base64.StdEncoding.EncodeToString(v)
	}

	return toYAML(m)
}

func (f files) Lines(path string) []string {
	if f == nil || f[path] == nil {
		return []string{}
	}

	return strings.Split(string(f[path]), "\n")
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
