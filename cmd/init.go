package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize conflogt in your repository",
	Long: `Initialize conflogt in your repository.

If no repository provided, this command will init conflogt in the current directory:
	1. create .conflogt directory inside your repository
	2. create .conflogt.yaml config file inside .conflogt
All files related to conflogt will be inside .conflogt`,
	Run: initRun,
}

func initRun(_ *cobra.Command, _ []string) {
	if viper.ConfigFileUsed() != "" {
		fmt.Println(fmt.Sprintf("conflogt is already initialized in `%s`", repository))
		os.Exit(1)
	}
	err := os.MkdirAll(cfgDir(), os.ModePerm)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create `%s` directory", cfgDir()))
		os.Exit(1)
	}
	err = ioutil.WriteFile(cfgFile(), []byte(""), 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("Unable to create `%s`", cfgFile()))
		os.Exit(1)
	}
	initConfig()
	fmt.Println(fmt.Sprintf("Successfully initialized conflogt in `%s`", repository))
	fmt.Println(fmt.Sprintf("A conflogt config file was created (`%s`)", cfgFile()))
	fmt.Println("Edit it according to your needs.")
}

func init() {
	rootCmd.AddCommand(initCmd)
}
