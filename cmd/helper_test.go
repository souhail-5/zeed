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
	writeCfgFile(t, []byte(""))
	writeChangelogFile(t, []byte(""))
}

func writeCfgFile(t *testing.T, data []byte) {
	if err := ioutil.WriteFile(cfgFile(), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func writeChangelogFile(t *testing.T, data []byte) {
	if err := ioutil.WriteFile(changelogFile(), data, 0644); err != nil {
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
	initCmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
	unifyCmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		flag.Value.Set(flag.DefValue)
	})
}
