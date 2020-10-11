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
		t.Fatalf("Expected %q got %q", 1, len(files))
	}
}
