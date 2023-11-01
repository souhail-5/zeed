package cmd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/souhail-5/zeed/internal/changelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"time"
)

var (
	channel string
	weight  int
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a zeed entry",
	Long:  `Add a zeed entry in your zeed repository.`,
	RunE:  addRun,
}

func addRun(_ *cobra.Command, args []string) error {
	entry := changelog.Entry{
		FrontMatter: changelog.FrontMatter{
			Channel: channel,
			Weight:  weight,
		},
		Text: args[0],
	}

	if _, err := entry.Validate(viper.GetViper()); err != nil {
		return errors.New(fmt.Sprintf("provided channel (\"%s\") is not supported", entry.FrontMatter.Channel))
	}

	return save(&entry)
}

func save(entry *changelog.Entry) error {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return err
	}
	filePath := filepath.Join(repository, ".zeed", id.String())
	yml, err := yaml.Marshal(&entry.FrontMatter)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(fmt.Sprintf("---\n%s---\n%s", yml, entry.Text)), 0644)
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&channel, "channel", "c", "default", "entry's channel")
	addCmd.Flags().IntVarP(&weight, "weight", "w", 0, "entry's weight")
}
