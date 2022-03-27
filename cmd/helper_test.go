package cmd

import (
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"testing"
)

func initRepository(t *testing.T) {
	var err error
	repository, err = ioutil.TempDir("", "repository")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(cfgDir(), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(cfgFile(), []byte(""), 0644)
	if err != nil {
		t.Fatal(err)
	}
}

func removeRepository(t *testing.T) {
	if err := os.RemoveAll(repository); err != nil {
		t.Fatal(err)
	}
}

func resetFlags() {
	rootCmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
}
