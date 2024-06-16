package cmd

import (
	"testing"
)

func TestNotInitializedProject(t *testing.T) {
	initRepository(t)
	rootCmd.SetArgs([]string{""}) // We need to reset args between each test
	expected := "zeed needs to be initialized in your repository. See `zeed init --help` for help"
	if err := rootCmd.Execute(); err != nil {
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	} else {
		t.Fatalf("Expected %q got %v", expected, nil)
	}
}

func TestInvalidCfgFile(t *testing.T) {
	initRepository(t)
	writeCfgFile(t, []byte(""))
	writeChangelogFile(t, []byte(""))
	defer removeRepository(t)
	writeCfgFile(t, []byte("invalid\nyaml"))
	rootCmd.SetArgs([]string{"add", "-t", "My changelog entry"})
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
	rootCmd.SetArgs([]string{"add", "-t", "My changelog entry"})
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
	writeCfgFile(t, []byte(""))
	writeChangelogFile(t, []byte(""))
	defer removeRepository(t)
	writeCfgFile(t, []byte("channels: ['feature', 'bugfix']"))
	rootCmd.SetArgs([]string{"add", "-t", "My changelog entry", "-c", "support"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested channel name must be considered not supported.")
	}
	expected := "entry's channel must be part of configured channels; channel (\"support\") is not supported"
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestInvalidTemplateFormat(t *testing.T) {
	resetFlags()
	initRepository(t)
	writeCfgFile(t, []byte(""))
	writeChangelogFile(t, []byte(""))
	defer removeRepository(t)
	writeCfgFile(t, []byte(`templates:
  bad-f0rmAt: "{{range .Entries}}- {{.Text}}\n{{end}}"
`))
	rootCmd.SetArgs([]string{"add", "-t", "My changelog entry"})
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
