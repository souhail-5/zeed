package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestAlreadyInitializedProject(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	rootCmd.SetArgs([]string{"init", repository})
	if err := rootCmd.Execute(); err != nil {
		expected := fmt.Sprintf("zeed is already initialized in `%s`", repository)
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}

	_ = os.Remove(cfgFile())
	_ = os.Chmod(cfgDir(), 0444)
	if err := rootCmd.Execute(); err != nil {
		expected := fmt.Sprintf("unable to create `%s`", cfgFile())
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}

	_ = os.Remove(cfgDir())
	_ = os.Chmod(repository, 0444)
	defer func() {
		_ = os.Chmod(repository, 0777)
	}()
	if err := rootCmd.Execute(); err != nil {
		expected := fmt.Sprintf("unable to create `%s` directory", cfgDir())
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}
}

func TestAlreadyInitializedProjectWithInvalidChannelFormat(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	err := ioutil.WriteFile(cfgFile(), []byte("channels:\n  - bad-f0rmAt"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"init", repository})
	if err = rootCmd.Execute(); err == nil {
		t.Fatal("An error must occurs.")
	}
	expected := fmt.Sprintf("zeed is already initialized in `%s`\ninvalid channel name: \"bad-f0rmAt\" (only a-z and _ are allowed)", repository)
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestInitialization(t *testing.T) {
	defer removeRepository(t)
	var err error
	repository, err := ioutil.TempDir("", "repository")
	if err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"init", "--repository", repository})
	if err := rootCmd.Execute(); err != nil {
		outputs := []string{
			fmt.Sprintf("Successfully initialized zeed in `%s`", repository),
			fmt.Sprintf("A zeed config file was created (`%s`)", cfgFile()),
			"Edit it according to your needs.",
		}
		expected := strings.Join(outputs, "\n")
		if expected != err.Error() {
			t.Fatalf("Expected %q got %q", expected, err.Error())
		}
	}
}
