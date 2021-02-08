package cmd

import (
	"io/ioutil"
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
	err := ioutil.WriteFile(cfgFile(), []byte("invalid\nyaml"), 0644)
	if err != nil {
		t.Fatal(err)
	}
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
	err := ioutil.WriteFile(cfgFile(), []byte("channels:\n  - bad-f0rmAt"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry"})
	if err = rootCmd.Execute(); err == nil {
		t.Fatal("Tested channel name must be considered invalid.")
	}
	expected := "Invalid channel name: \"bad-f0rmAt\". Only a-z and _ are allowed."
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestEntry(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	rootCmd.SetArgs([]string{"My changelog entry"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	files, err := entriesFiles()
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("Repository must have 1 entry file, got %d (%v)", len(files), files)
	}
	file := files[0]
	if file.Channel != "undefined" {
		t.Errorf("Default channel must be %q, got %q", "undefined", file.Channel)
	}
	if file.Priority != 0 {
		t.Errorf("Default priority must be %q, got %q", "0", file.Priority)
	}
	expectedContent := "My changelog entry"
	if file.Content != expectedContent {
		t.Errorf("Entry file content must be %q, got %q", expectedContent, file.Content)
	}
}
