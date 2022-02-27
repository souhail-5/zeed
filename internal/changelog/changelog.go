package changelog

import (
	"embed"
	"errors"

	"github.com/spf13/viper"
)

//go:embed template
var Templates embed.FS

type File struct {
	Name     string
	Channel  string
	Priority int
	Hash     string
	Content  string
}

type Entry struct {
	Text     string
	Priority int
	Channel  Channel
}

type Channel struct {
	Id      string
	Entries []Entry
}

type ByPriority []Entry

func (entries ByPriority) Len() int           { return len(entries) }
func (entries ByPriority) Swap(i, j int)      { entries[i], entries[j] = entries[j], entries[i] }
func (entries ByPriority) Less(i, j int) bool { return entries[i].Priority > entries[j].Priority }

func (file File) Validate(viper *viper.Viper) (ok bool, err error) {
	channels := viper.GetStringSlice("channels")
	if len(channels) != 0 && !Contains(channels, file.Channel) {
		return false, errors.New("channel's file must be part of configured channels")
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
