package document

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Near identical to https://github.com/helm/helm/blob/main/pkg/engine/files.go as to preserve the interface.

type fileEntry struct {
	Path string
	data []byte
}

func (f *fileEntry) GetData() []byte {
	if f.data == nil {
		data, err := ioutil.ReadFile(f.Path)
		if err != nil {
			log.Warnf("Error reading file contents for %s: %s", f.Path, err.Error())
			return []byte{}
		}
		f.data = data
	}

	return f.data
}

type files map[string]*fileEntry

func getFiles(dir string) (files, error) {
	result := make(files)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		result[path] = &fileEntry{Path: path}
		return nil
	})
	if err != nil {
		return map[string]*fileEntry{}, err
	}

	return result, nil
}

func (f files) GetBytes(name string) []byte {
	if v, ok := f[name]; ok {
		return v.GetData()
	}
	return []byte{}
}

func (f files) Get(name string) string {
	return string(f.GetBytes(name))
}

func (f files) Glob(pattern string) files {
	result := make(files)
	g, err := glob.Compile(pattern, '/')
	if err != nil {
		log.Warnf("Error compiling Glob patten %s: %s", pattern, err.Error())
		return result
	}

	for filePath, entry := range f {
		if g.Match(filePath) {
			result[filePath] = entry
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
		m[path.Base(k)] = string(v.GetData())
	}

	return toYAML(m)
}

func (f files) AsSecrets() string {
	if f == nil {
		return ""
	}

	m := make(map[string]string)

	for k, v := range f {
		m[path.Base(k)] = base64.StdEncoding.EncodeToString(v.GetData())
	}

	return toYAML(m)
}

func (f files) Lines(path string) []string {
	if f == nil {
		return []string{}
	}
	entry, exists := f[path]
	if !exists {
		return []string{}
	}

	return strings.Split(string(entry.GetData()), "\n")
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}
