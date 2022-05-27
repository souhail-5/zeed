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
	writeCfgFile(t, []byte("channels:\n  - bad-f0rmAt"))
	rootCmd.SetArgs([]string{"init", repository})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested initialization should fail because an invalid channel name is in the config file.")
	}
	expected := fmt.Sprintf("zeed is already initialized in `%s`\ninvalid channel name: \"bad-f0rmAt\" (only a-z and _ are allowed)", repository)
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestAlreadyInitializedProjectWithInvalidTemplateFormat(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	writeCfgFile(t, []byte(`templates:
  bad-f0rmAt: "{{range .Entries}}- {{.Text}}\n{{end}}"
`))
	rootCmd.SetArgs([]string{"init", repository})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested initialization should fail because an invalid template name is in the config file.")
	}
	// TODO https://github.com/souhail-5/zeed/issues/16
	expected := fmt.Sprintf("zeed is already initialized in `%s`\ninvalid template name: \"bad-f0rmat\" (only a-z and _ are allowed)", repository)
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
