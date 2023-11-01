package cmd

import (
	"crypto/rand"
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
	text    string
	weight  int
)

var addCmd = &cobra.Command{
	Use:     "add -t text [-c channel] [-w weight]",
	Example: "zeed add -t \"Add zeed config to the repository.\" -c added -w 128",
	Short:   "Add a zeed entry",
	Long:    `Add a zeed entry in your zeed repository.`,
	RunE:    addRun,
}

func addRun(_ *cobra.Command, args []string) error {
	entry := changelog.Entry{
		FrontMatter: changelog.FrontMatter{
			Channel: channel,
			Weight:  weight,
		},
		Text: text,
	}

	if _, err := entry.Validate(viper.GetViper()); err != nil {
		return err
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
	addCmd.Flags().StringVarP(&text, "text", "t", "", "entry's text")
	addCmd.Flags().StringVarP(&channel, "channel", "c", "default", "entry's channel")
	addCmd.Flags().IntVarP(&weight, "weight", "w", 0, "entry's weight")
}
