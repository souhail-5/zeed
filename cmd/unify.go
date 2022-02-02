package cmd

import (
	"errors"
	"fmt"
	"github.com/souhail-5/zeed/internal/changelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var unifyCmd = &cobra.Command{
	Use:   "unify",
	Short: "Print unified changelog entries",
	Long:  `Print unified changelog entries.`,
	RunE:  unifyRun,
}

func init() {
	rootCmd.AddCommand(unifyCmd)
	unifyCmd.Flags().Bool("flush", false, "if set, entries will be removed after `unify`")
	unifyCmd.Flags().StringP("template", "t", "", "unify template")
	if err := viper.BindPFlag("template", unifyCmd.Flags().Lookup("template")); err != nil {
		rootCmd.Println("Unable to get template config.")
		os.Exit(1)
	}
}

func unifyRun(cmd *cobra.Command, _ []string) error {
	files, err := entriesFiles()
	if err != nil {
		return errors.New("unable to read zeed files")
	}
	var data struct {
		Entries  []changelog.Entry
		Channels map[string]changelog.Channel
	}
	data.Entries, data.Channels = entries(files)

	var tmpl *template.Template
	if path := viper.GetString("template"); path != "" {
		tmpl = template.New(path)
		tmpl, err = tmpl.ParseFiles(filepath.Join(cfgDir(), path))
	} else {
		tmpl = template.New("default")
		tmpl, err = tmpl.Parse("{{range .Entries}}{{.Text}}\n{{end}}")
	}
	if err != nil {
		return errors.New("unable to read zeed template")
	}
	err = tmpl.Execute(cmd.OutOrStdout(), data)
	if err != nil {
		return errors.New("unable to unify")
	}
	if shouldFlush, _ := cmd.Flags().GetBool("flush"); shouldFlush {
		for _, file := range files {
			os.Remove(filepath.Join(cfgDir(), file.Name))
		}
	}

	return nil
}

func entries(files []changelog.File) ([]changelog.Entry, map[string]changelog.Channel) {
	var entries []changelog.Entry
	var channels map[string]changelog.Channel
	channels = make(map[string]changelog.Channel)

	for _, file := range files {
		if _, exist := channels[file.Channel]; !exist {
			channels[file.Channel] = changelog.Channel{
				Id: file.Channel,
			}
		}
		entry := changelog.Entry{
			Text:     file.Content,
			Priority: file.Priority,
			Channel:  channels[file.Channel],
		}
		entries = append(entries, entry)
		channel := channels[file.Channel]
		channel.Entries = append(channel.Entries, entry)
		channels[file.Channel] = channel
	}
	sort.Sort(changelog.ByPriority(entries))
	for _, channel := range channels {
		sort.Slice(channel.Entries, func(i, j int) bool {
			return channel.Entries[i].Priority > channel.Entries[j].Priority
		})
	}

	return entries, channels
}

func entriesFiles() ([]changelog.File, error) {
	var files []changelog.File
	f, err := os.Open(cfgDir())
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	err = f.Close()
	if err != nil {
		return files, err
	}
	channels := viper.GetStringSlice("channels")

	for _, file := range fileInfo {
		if file.Name() == filepath.Base(cfgFile()) {
			continue
		}
		metadata := strings.Split(file.Name(), "=")
		if len(metadata) != 3 {
			continue
		}
		channel := metadata[0]
		if !contains(channels, channel) && channel != "undefined" {
			fmt.Println("Entry \"" + file.Name() + "\" was not processed: its channel is not supported")
			continue
		}
		// TODO verify file extension
		// TODO ignore template files
		content, _ := ioutil.ReadFile(filepath.Join(cfgDir(), file.Name()))
		priority, _ := strconv.Atoi(metadata[1])
		files = append(files, changelog.File{
			Name:     file.Name(),
			Channel:  channel,
			Priority: priority,
			Hash:     strings.Split(metadata[2], ".")[0],
			Content:  string(content),
		})
	}

	return files, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
