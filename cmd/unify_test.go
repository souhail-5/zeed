package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestUnify(t *testing.T) {
	initRepository(t)
	defer removeRepository(t)
	bufferString := bytes.NewBufferString("")
	rootCmd.SetOut(bufferString)
	rootCmd.SetArgs([]string{"My changelog entry #1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #2", "-p", fmt.Sprintf("%d", 32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"unify"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(bufferString)
	if err != nil {
		t.Fatal(err)
	}
	expected := "My changelog entry #2\nMy changelog entry #1\n"
	if expected != string(out) {
		t.Fatalf("Expected %q got %q", expected, string(out))
	}

	rootCmd.SetArgs([]string{"unify", "--flush"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	files, err := ioutil.ReadDir(cfgDir())
	if len(files) != 1 {
		if err != nil {
			t.Fatal(err)
		}
		t.Fatalf("Expected %v got %v", 1, len(files))
	}
}

func TestUnifyWithTemplate(t *testing.T) {
	resetFlags()
	initRepository(t)
	writeCfgFile(t, []byte("channels: [added, security]"))
	defer removeRepository(t)
	bufferString := bytes.NewBufferString("")
	rootCmd.SetOut(bufferString)
	rootCmd.SetArgs([]string{"My changelog entry #1", "-c", "added"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{
		"My changelog entry #2",
		"-p",
		fmt.Sprintf("%d", 32),
		"-c",
		"added",
	})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #3", "-c", "security"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"unify", "-t", "keepachangelog"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(bufferString)
	if err != nil {
		t.Fatal(err)
	}
	expected := `### Added
- My changelog entry #2
- My changelog entry #1

### Security
- My changelog entry #3

`
	if expected != string(out) {
		t.Fatalf("Expected %q got %q", expected, string(out))
	}
}

func TestUnifyWithConfiguredTemplate(t *testing.T) {
	resetFlags()
	initRepository(t)
	writeCfgFile(t, []byte(`templates:
  default: "{{range .Entries}}• {{.Text}}\n{{end}}"
`))
	defer removeRepository(t)
	bufferString := bytes.NewBufferString("")
	rootCmd.SetOut(bufferString)
	rootCmd.SetArgs([]string{"My changelog entry #1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #2", "-p", fmt.Sprintf("%d", 32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #3", "-p", fmt.Sprintf("%d", 16)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"unify"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(bufferString)
	if err != nil {
		t.Fatal(err)
	}
	expected := `• My changelog entry #2
• My changelog entry #3
• My changelog entry #1
`
	if expected != string(out) {
		t.Fatalf("Expected %q got %q", expected, string(out))
	}
}
