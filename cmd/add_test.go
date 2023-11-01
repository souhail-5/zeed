package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAdd(t *testing.T) {
	resetFlags()
	initRepository(t)
	writeCfgFile(t, []byte(""))
	writeChangelogFile(t, []byte(""))
	defer removeRepository(t)
	rootCmd.SetArgs([]string{"add", "-t", "My changelog entry"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(cfgDir())
	if err != nil {
		t.Fatal(err)
	}
	fileInfo, err := f.Readdir(-1)
	err = f.Close()
	if err != nil {
		t.Fatal(err)
	}
	expectedContent := `---
channel: default
weight: 0
---
My changelog entry`
	for _, file := range fileInfo {
		if file.Name() == filepath.Base(cfgFile()) {
			continue
		}
		content, _ := ioutil.ReadFile(filepath.Join(cfgDir(), file.Name()))
		if string(content) != expectedContent {
			t.Errorf("entry file content must be %q, got %q", expectedContent, content)
		}
	}
}
