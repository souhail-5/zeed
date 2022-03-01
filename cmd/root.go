package cmd

import (
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/souhail-5/zeed/internal/changelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// const improve performance
const ALPH = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var (
	isCfgFileLoaded bool
	verbose         bool
	repository      string
	cchannel        string
	priority        int
)

var rootCmd = &cobra.Command{
	Use:     "zeed <entry_text>",
	Example: "zeed \"Add zeed config to the repository.\" -c added -p 128",
	Version: "1.0.0-beta",
	Short:   "A tool to eliminate changelog-related merge conflicts",
	Long: `Zeed is a free and open source tool
to eliminate changelog-related merge conflicts.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.SilenceUsage = false
			return fmt.Errorf("accepts %d arg(s), received %d", 1, len(args))
		}
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Use != "init" {
			if viper.ConfigFileUsed() == "" {
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
	RunE:         rootRun,
	SilenceUsage: true,
}

func rootRun(_ *cobra.Command, args []string) error {
	file := changelog.File{
		Channel:  cchannel,
		Priority: priority,
		Content:  args[0],
	}

	if _, err := file.Validate(viper.GetViper()); err != nil {
		return errors.New(fmt.Sprintf("provided channel (\"%s\") is not supported", file.Channel))
	}

	return save(&file)
}

func save(file *changelog.File) error {
	id, err := gonanoid.Generate(ALPH, 10)
	if err != nil {
		return err
	}
	file.Hash = id
	file.Name = strings.Join([]string{file.Channel, strconv.Itoa(file.Priority), file.Hash}, "=")
	filePath := filepath.Join(repository, ".zeed", file.Name)

	return ioutil.WriteFile(filePath, []byte(file.Content), 0644)
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVar(&repository, "repository", "", "path to your project's repository")
	rootCmd.Flags().StringVarP(&cchannel, "channel", "c", "default", "entry's channel")
	rootCmd.Flags().IntVarP(&priority, "priority", "p", 0, "entry's priority")
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
	for _, channel := range viper.GetStringSlice("channels") {
		if match, _ := regexp.MatchString("^[a-z_]+$", channel); !match {
			return errors.New(fmt.Sprintf("invalid channel name: \"%s\" (only a-z and _ are allowed)", channel))
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
