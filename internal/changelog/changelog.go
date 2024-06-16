package changelog

import (
	"embed"
	"errors"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/spf13/viper"
	"io"
)

//go:embed template
var Templates embed.FS

type FrontMatter struct {
	Channel string `yaml:"channel,omitempty"`
	Weight  int    `yaml:"weight"`
}

type Entry struct {
	FrontMatter FrontMatter
	Text        string
}

type Channel struct {
	Id      string
	Entries []*Entry
}

func NewEntry(r io.Reader) (*Entry, error) {
	fm := FrontMatter{}

	rest, err := frontmatter.MustParse(r, &fm)
	if err != nil {
		return nil, err
	}

	return &Entry{
		FrontMatter: fm,
		Text:        string(rest),
	}, nil
}

type Entries []*Entry

func (entries Entries) Len() int      { return len(entries) }
func (entries Entries) Swap(i, j int) { entries[i], entries[j] = entries[j], entries[i] }
func (entries Entries) Less(i, j int) bool {
	return entries[i].FrontMatter.Weight > entries[j].FrontMatter.Weight
}

func (e Entry) Validate(viper *viper.Viper) (ok bool, err error) {
	if e.Text == "" {
		return false, errors.New("entry's text must not be empty")
	}

	channels := viper.GetStringSlice("channels")
	if e.FrontMatter.Channel != "" && len(channels) != 0 && !Contains(channels, e.FrontMatter.Channel) {
		return false, errors.New(fmt.Sprintf("entry's channel must be part of configured channels; channel (\"%s\") is not supported", e.FrontMatter.Channel))
	}

	return true, nil
}

// Contains https://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func Contains(s []string, i string) bool {
	for _, e := range s {
		if e == i {
			return true
		}
	}

	return false
}
