package cmd

import (
	"errors"
	"fmt"
	"github.com/souhail-5/zeed/internal/changelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"path/filepath"
	"sort"
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
	files, err := files()
	if err != nil {
		return errors.New("unable to read zeed files")
	}
	var data struct {
		Entries  []*changelog.Entry
		Channels map[string]changelog.Channel
	}
	data.Entries, data.Channels = entries(files)

	var tmpl *template.Template
	if t := viper.GetString("template"); t != "" {
		configFS := os.DirFS(cfgDir())
		tmpl = template.New(t)
		tmpl, err = tmpl.ParseFS(configFS, t)
		if err != nil {
			tmpl, err = tmpl.ParseFS(changelog.Templates, filepath.Join("template", t))
		}
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
			os.Remove(file.Name())
		}
	}

	return nil
}

func files() ([]*os.File, error) {
	var files []*os.File
	d, err := os.Open(cfgDir())
	if err != nil {
		return files, err
	}
	fileInfos, err := d.Readdir(-1)
	err = d.Close()
	if err != nil {
		return files, err
	}

	for _, info := range fileInfos {
		if info.Name() == filepath.Base(cfgFile()) {
			continue
		}
		file, err := os.Open(filepath.Join(cfgDir(), info.Name()))
		if err != nil {
			return []*os.File{}, err
		}
		files = append(files, file)
	}

	return files, nil
}

func entries(files []*os.File) ([]*changelog.Entry, map[string]changelog.Channel) {
	var entries []*changelog.Entry
	var channels map[string]changelog.Channel
	channels = make(map[string]changelog.Channel)
	cc := viper.GetStringSlice("channels")

	for _, file := range files {
		e, err := changelog.NewEntry(file)
		if err != nil {
			fmt.Println(err.Error())
		}
		if !changelog.Contains(cc, e.FrontMatter.Channel) && e.FrontMatter.Channel != "default" {
			fmt.Println("entry \"" + file.Name() + "\" was not processed: its channel is not supported")
			continue
		}
		if _, exist := channels[e.FrontMatter.Channel]; !exist {
			channels[e.FrontMatter.Channel] = changelog.Channel{
				Id: e.FrontMatter.Channel,
			}
		}
		entries = append(entries, e)
		channel := channels[e.FrontMatter.Channel]
		channel.Entries = append(channel.Entries, e)
		channels[e.FrontMatter.Channel] = channel
	}
	sort.Sort(changelog.ByPriority(entries))
	for _, channel := range channels {
		sort.Slice(channel.Entries, func(i, j int) bool {
			return channel.Entries[i].FrontMatter.Priority > channel.Entries[j].FrontMatter.Priority
		})
	}

	return entries, channels
}
