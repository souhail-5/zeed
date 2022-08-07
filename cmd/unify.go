package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/souhail-5/zeed/internal/changelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	aline string
	bline string
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
	unifyCmd.Flags().StringP("template", "t", "default", "unify template")
	unifyCmd.Flags().StringVarP(&aline, "aline", "a", "", "the line after which the unified entries will be pasted")
	unifyCmd.Flags().StringVarP(&bline, "bline", "b", "", "the line before which the unified entries will be pasted")
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

	k, err := cmd.Flags().GetString("template")
	if err != nil {
		return errors.New("unable to read template flag")
	}
	s := viper.GetString("templates." + k)
	tmpl := template.New(k)
	tmpl, err = tmpl.Parse(s)
	if err != nil || s == "" {
		tmpl, err = tmpl.ParseFS(changelog.Templates, filepath.Join("template", k))
		if err != nil {
			return errors.New(fmt.Sprintf("provided template (\"%s\") is not supported", k))
		}
	}
	unifiedText := bytes.Buffer{}
	mw := io.MultiWriter(&unifiedText, cmd.OutOrStdout())
	err = tmpl.Execute(mw, data)
	if err != nil {
		return errors.New("unable to unify")
	}

	err = uChangelog(filepath.Join(repository, "CHANGELOG.md"), unifiedText.String(), aline, bline)
	if err != nil {
		return err
	}

	if shouldFlush, _ := cmd.Flags().GetBool("flush"); shouldFlush {
		for _, file := range files {
			err = os.Remove(file.Name())
		}
		if err != nil {
			return errors.New("unable to remove all the entries")
		}
	}

	return nil
}

func uChangelog(filename string, unifiedText string, aline string, bline string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("unable to read the changelog file")
	}

	var re string
	var repl []byte

	if aline != "" && bline == "" {
		re = fmt.Sprintf(`(?s)^(.*?\n?)(?-s)(.*%s.*\n?)(?s)(.*)`, aline)
		repl = []byte(fmt.Sprintf(`${1}${2}%s${3}`, unifiedText))
	}

	if bline != "" && aline == "" {
		re = fmt.Sprintf(`(?s)^(.*?\n?)(?-s)(.*%s.*\n?)(?s)(.*)`, bline)
		repl = []byte(fmt.Sprintf(`${1}%s${2}${3}`, unifiedText))
	}

	if aline != "" && bline != "" {
		re = fmt.Sprintf(`(?s)^(.*?\n?)(?-s)(.*%s.*\n?)(?s)(.*?\n?)(?-s)(.*%s.*\n?)(?s)(.*)`, aline, bline)
		repl = []byte(fmt.Sprintf(`${1}${2}%s${4}${5}`, unifiedText))
	}

	if aline != "" || bline != "" {
		r, err := regexp.Compile(re)
		if err != nil {
			return errors.New("unable to compile the regex with aline flag")
		}
		content = r.ReplaceAll(content, repl)
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			return errors.New("unable to write the changelog file")
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
	sort.Sort(changelog.ByWeight(entries))
	for _, channel := range channels {
		sort.Slice(channel.Entries, func(i, j int) bool {
			return channel.Entries[i].FrontMatter.Weight > channel.Entries[j].FrontMatter.Weight
		})
	}

	return entries, channels
}
