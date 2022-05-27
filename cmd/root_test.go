package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNotInitializedProject(t *testing.T) {
	rootCmd.SetArgs([]string{"My changelog entry"})
	if err := rootCmd.Execute(); err != nil {
		expected := "zeed needs to be initialized in your repository. See `zeed init --help` for help"
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}
}

func TestInvalidCfgFile(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	writeCfgFile(t, []byte("invalid\nyaml"))
	rootCmd.SetArgs([]string{"My changelog entry"})
	if err := rootCmd.Execute(); err != nil {
		expected := "unable to read your config file"
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}
}

func TestInvalidChannelFormat(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	writeCfgFile(t, []byte("channels:\n  - bad-f0rmAt"))
	rootCmd.SetArgs([]string{"My changelog entry"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested channel name must be considered invalid.")
	}
	expected := "invalid channel name: \"bad-f0rmAt\" (only a-z and _ are allowed)"
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestUnconfiguredChannel(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	writeCfgFile(t, []byte("channels: ['feature', 'bugfix']"))
	rootCmd.SetArgs([]string{"My changelog entry", "-c", "support"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested channel name must be considered not supported.")
	}
	expected := "provided channel (\"support\") is not supported"
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestInvalidTemplateFormat(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	writeCfgFile(t, []byte(`templates:
  bad-f0rmAt: "{{range .Entries}}- {{.Text}}\n{{end}}"
`))
	rootCmd.SetArgs([]string{"My changelog entry"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested template name must be considered invalid.")
	}
	// TODO https://github.com/souhail-5/zeed/issues/16
	expected := "invalid template name: \"bad-f0rmat\" (only a-z and _ are allowed)"
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestRoot(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	rootCmd.SetArgs([]string{"My changelog entry"})
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
