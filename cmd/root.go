package cmd

import (
	"fmt"
	"github.com/Souhail-5/conflogt/changelog"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// const improve performance
const ALPH = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	isCfgFileLoaded bool
	repository string
	cchannel string
	priority int
)

var rootCmd = &cobra.Command{
	Use:   "conflogt",
	Version: "1.0.0",
	Short: "A tool to eliminate changelog-related merge conflicts",
	Long: `Conflogt is a free and open source tool
to eliminate changelog-related merge conflicts.`,
	Args: cobra.ExactArgs(1),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Use != "init" {
			if viper.ConfigFileUsed() == "" {
				fmt.Println("conflogt needs to be initialized in your repository.", "Did you forget to use `conflogt init`? (see `conflogt init --help` for more help and options)")
				os.Exit(1)
			} else if !isCfgFileLoaded {
				fmt.Println("Unable to read your config file.")
				os.Exit(1)
			}
		}
	},
	Run: rootRun,
}

func rootRun(_ *cobra.Command, args []string) {
	file := changelog.File{
		Channel:  cchannel,
		Priority: priority,
		Content:  args[0],
	}
	err := save(&file)
	if err != nil {
		fmt.Println(err)
	}
}

func save(file *changelog.File) error {
	id, err := gonanoid.Generate(ALPH, 10)
	if err != nil {
		return err
	}
	file.Hash = id
	file.Name = strings.Join([]string{file.Channel, strconv.Itoa(file.Priority), file.Hash}, "=")
	filePath := filepath.Join(repository, ".conflogt", file.Name)

	return ioutil.WriteFile(filePath, []byte(file.Content), 0644)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "path to your project's repository")
	rootCmd.Flags().StringVarP(&cchannel, "channel", "c", "undefined", "Entry's channel")
	rootCmd.Flags().IntVarP(&priority, "priority", "p", 0, "Entry's priority")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if repository != "" {
		viper.AddConfigPath(cfgDir())
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		repository = wd
		searchPath := wd
		for ; searchPath != string(os.PathSeparator) ; searchPath = filepath.Dir(searchPath) {
			viper.AddConfigPath(filepath.Join(searchPath, ".conflogt"))
		}
		viper.AddConfigPath(filepath.Join(searchPath, ".conflogt"))
	}
	viper.SetConfigName(".conflogt")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("conflogt")
	if err := viper.ReadInConfig(); err == nil {
		isCfgFileLoaded = true
		repository = filepath.Dir(filepath.Dir(viper.ConfigFileUsed()))
		fmt.Println("Running for:", repository)
	}
}

func cfgDir() string {
	return filepath.Join(repository, ".conflogt")
}

func cfgFile() string {
	return filepath.Join(cfgDir(), ".conflogt.yaml")
}
