package changelog

import (
	"embed"
	"errors"
	"github.com/adrg/frontmatter"
	"github.com/spf13/viper"
	"io"
)

//go:embed template
var Templates embed.FS

type FrontMatter struct {
	Channel  string `yaml:"channel"`
	Priority int    `yaml:"priority"`
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

type ByPriority []*Entry

func (entries ByPriority) Len() int      { return len(entries) }
func (entries ByPriority) Swap(i, j int) { entries[i], entries[j] = entries[j], entries[i] }
func (entries ByPriority) Less(i, j int) bool {
	return entries[i].FrontMatter.Priority > entries[j].FrontMatter.Priority
}

func (e Entry) Validate(viper *viper.Viper) (ok bool, err error) {
	channels := viper.GetStringSlice("channels")
	if len(channels) != 0 && !Contains(channels, e.FrontMatter.Channel) {
		return false, errors.New("channel's entry must be part of configured channels")
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
