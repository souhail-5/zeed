package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
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
	rootCmd.SetArgs([]string{"My changelog entry #2", "-w", fmt.Sprintf("%d", 32)})
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
	defer removeRepository(t)
	writeCfgFile(t, []byte("channels: [added, security]"))
	bufferString := bytes.NewBufferString("")
	rootCmd.SetOut(bufferString)
	rootCmd.SetArgs([]string{"My changelog entry #1", "-c", "added"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{
		"My changelog entry #2",
		"-w",
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
	defer removeRepository(t)
	writeCfgFile(t, []byte(`templates:
  default: "{{range .Entries}}• {{.Text}}\n{{end}}"
`))
	bufferString := bytes.NewBufferString("")
	rootCmd.SetOut(bufferString)
	rootCmd.SetArgs([]string{"My changelog entry #1"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #2", "-w", fmt.Sprintf("%d", 32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	rootCmd.SetArgs([]string{"My changelog entry #3", "-w", fmt.Sprintf("%d", 16)})
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

func TestUnifyWithUnconfiguredTemplate(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	rootCmd.SetArgs([]string{"unify", "-t", "slack"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("Tested template name must be considered not supported.")
	}
	expected := "provided template (\"slack\") is not supported"
	if expected != err.Error() {
		t.Fatalf("Expected %q got %q", expected, err.Error())
	}
}

func TestUnifyAline(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	writeChangelogFile(t, []byte(`# Changelog
A short introduction

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`))
	rootCmd.SetArgs([]string{"--", "- lorem ipsum"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"## Version Tortank", "-w", fmt.Sprintf("%d", 32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"", "-w", fmt.Sprintf("%d", 64)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"unify", "-a", "A short introduction"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	expected := `# Changelog
A short introduction

## Version Tortank
- lorem ipsum

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`
	c, err := ioutil.ReadFile(filepath.Join(repository, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(c) {
		t.Fatalf("Expected %q got %q", expected, c)
	}
}

func TestUnifyBline(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	writeChangelogFile(t, []byte(`# Changelog
A short introduction

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`))
	rootCmd.SetArgs([]string{"--", "- lorem ipsum"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"## Version Tortank", "-w", fmt.Sprintf("%d", 32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"", "-w", fmt.Sprintf("%d", -32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"unify", "-b", "## Version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	expected := `# Changelog
A short introduction

## Version Tortank
- lorem ipsum

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`
	c, err := ioutil.ReadFile(filepath.Join(repository, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(c) {
		t.Fatalf("Expected %q got %q", expected, c)
	}
}

func TestUnifyAlineBline(t *testing.T) {
	resetFlags()
	initRepository(t)
	defer removeRepository(t)
	writeChangelogFile(t, []byte(`# Changelog
A short introduction

## Unreleased
Nothing here.

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`))
	rootCmd.SetArgs([]string{"--", "- lorem ipsum"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"", "-w", fmt.Sprintf("%d", -32)})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"unify", "-a", "## Unreleased", "-b", "## Version"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}

	expected := `# Changelog
A short introduction

## Unreleased
- lorem ipsum

## Version Carabaffe
- lorem ipsum

## Version Carapuce
- lorem ipsum
- lorem ipsum
`
	c, err := ioutil.ReadFile(filepath.Join(repository, "CHANGELOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if expected != string(c) {
		t.Fatalf("Expected %q got %q", expected, c)
	}
}
