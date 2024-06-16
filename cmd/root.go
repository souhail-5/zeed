package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
)

var (
	cmdErrInitBus   = newErrInitBus()
	isCfgFileLoaded bool
	repository      string
	verbose         bool
)

var rootCmd = &cobra.Command{
	Use:   "zeed",
	Short: "A tool to eliminate changelog-related merge conflicts.",
	Long:  `Zeed is a free and open source tool to eliminate changelog-related merge conflicts.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmdErrInitBus.errors != nil {
			return cmdErrInitBus
		}
		if cmd.Use != "init" {
			if repository == "" || viper.ConfigFileUsed() == "" {
				return errors.New("zeed needs to be initialized in your repository. See `zeed init --help` for help")
			} else if !isCfgFileLoaded {
				return errors.New("unable to read your config file")
			} else if err := validateConfig(viper.GetViper()); err != nil {
				return err
			}
			if verbose {
				cmd.Println("Running for:", repository)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running for:", repository)
		return nil
	},
	SilenceUsage: true,
}

func SetVersion(v string) {
	rootCmd.Version = v
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	if cmdErrInitBus.errors != nil {
		return
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "path to your project's repository")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if repository != "" {
		viper.AddConfigPath(cfgDir())
	} else {
		wd, err := os.Getwd()
		if err != nil {
			cmdErrInitBus.AppendError(err)
			return
		}
		repository = wd
		searchPath := wd
		for ; searchPath != string(os.PathSeparator); searchPath = filepath.Dir(searchPath) {
			viper.AddConfigPath(filepath.Join(searchPath, ".zeed"))
		}
		viper.AddConfigPath(filepath.Join(searchPath, ".zeed"))
	}
	viper.SetConfigName(".zeed")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("zeed")
	if err := viper.ReadInConfig(); err == nil {
		isCfgFileLoaded = true
		repository = filepath.Dir(filepath.Dir(viper.ConfigFileUsed()))
	}
}

func validateConfig(viper *viper.Viper) error {
	for _, cn := range viper.GetStringSlice("channels") {
		if match, _ := regexp.MatchString("^[a-z_]+$", cn); !match {
			return errors.New(fmt.Sprintf("invalid channel name: \"%s\" (only a-z and _ are allowed)", cn))
		}
	}
	for tn := range viper.GetStringMap("templates") {
		if match, _ := regexp.MatchString("^[a-z_]+$", tn); !match {
			return errors.New(fmt.Sprintf("invalid template name: \"%s\" (only a-z and _ are allowed)", tn))
		}
	}

	return nil
}

func cfgDir() string {
	return filepath.Join(repository, ".zeed")
}

func cfgFile() string {
	return filepath.Join(cfgDir(), ".zeed.yaml")
}
