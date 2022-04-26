package document

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

// As the interface has been kept the same as in Helm, the tests also work here.
// Tests similar to https://github.com/helm/helm/blob/main/pkg/engine/files_test.go.

var cases = []struct {
	path, data string
}{
	{"ship/captain.txt", "The Captain"},
	{"ship/stowaway.txt", "Legatt"},
	{"story/name.txt", "The Secret Sharer"},
	{"story/author.txt", "Joseph Conrad"},
	{"multiline/test.txt", "bar\nfoo"},
}

func getTestFiles() files {
	a := make(files, len(cases))
	for _, c := range cases {
		a[c.path] = []byte(c.data)
	}
	return a
}

func TestNewFiles(t *testing.T) {
	files := getTestFiles()
	if len(files) != len(cases) {
		t.Errorf("Expected len() = %d, got %d", len(cases), len(files))
	}

	for i, f := range cases {
		if got := string(files.GetBytes(f.path)); got != f.data {
			t.Errorf("%d: expected %q, got %q", i, f.data, got)
		}
		if got := files.Get(f.path); got != f.data {
			t.Errorf("%d: expected %q, got %q", i, f.data, got)
		}
	}
}

func TestFileGlob(t *testing.T) {
	as := assert.New(t)

	f := getTestFiles()

	matched := f.Glob("story/**")

	as.Len(matched, 2, "Should be two files in glob story/**")
	as.Equal("Joseph Conrad", matched.Get("story/author.txt"))
}

func TestToConfig(t *testing.T) {
	as := assert.New(t)

	f := getTestFiles()
	out := f.Glob("**/captain.txt").AsConfig()
	as.Equal("captain.txt: The Captain", out)

	out = f.Glob("ship/**").AsConfig()
	as.Equal("captain.txt: The Captain\nstowaway.txt: Legatt", out)
}

func TestToSecret(t *testing.T) {
	as := assert.New(t)

	f := getTestFiles()

	out := f.Glob("ship/**").AsSecrets()
	as.Equal("captain.txt: VGhlIENhcHRhaW4=\nstowaway.txt: TGVnYXR0", out)
}

func TestLines(t *testing.T) {
	as := assert.New(t)

	f := getTestFiles()

	out := f.Lines("multiline/test.txt")
	as.Len(out, 2)

	as.Equal("bar", out[0])
}

func TestGetFiles(t *testing.T) {
	chartDir, err := os.MkdirTemp("", "*-helm-docs-chart")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(chartDir)
	})

	testFiles := getTestFiles()
	for filePath, fileData := range testFiles {
		fullPath := path.Join(chartDir, filePath)
		baseDir := path.Dir(fullPath)
		if err = os.MkdirAll(baseDir, 0o755); err != nil {
			t.Fatal(err)
		}

		if err = os.WriteFile(fullPath, fileData, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	chartFiles, err := getFiles(chartDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(chartFiles) != len(testFiles) {
		t.Errorf("chart files: expected %d, got %d", len(chartFiles), len(testFiles))
	}

	// Sanity check the files have been read
	for filePath, data := range chartFiles {
		if len(data) == 0 {
			t.Errorf("%s: expected file contents, got 0 bytes", filePath)
		}
	}
}
