# Zeed
Zeed is a free and open source tool to eliminate changelog-related merge conflicts. Team collaboration and continuous integration will be easier.

**How it works?** Use Zeed to add any entry to your changelog. Zeed will not modify your changelog file, but save your entries as a file within a staging area. When your are ready, you ask zeed to unify all staged entries, and render them according to a template. Copy/Paste the rendering to your changelog file. Use Zeed to delete the staged entries and start over for another release.

## Getting Started
These instructions will get you zeed up and running on your local machine.

### Technical prerequisites
Zeed is written with support for multiple platforms. Zeed currently provides binaries for the following:

- macOS (Darwin) for x64, i386, and ARM architectures
- Windows
- Linux

### How to install?

Download the appropriate version for your platform from [Zeed Releases](https://github.com/souhail-5/zeed/releases). Once downloaded, the binary can be run from anywhere. You don’t need to install it into a global location. This works well for shared hosts and other systems where you don’t have a privileged account.

Ideally, you should install it somewhere in your PATH for easy use. /usr/local/bin is the most probable location.

Verify your installation by running `zeed --version`

## Basic Usage

- Init zeed within your project `zeed init`
- Add an entry `zeed "I am a changelog entry"`
- Add another entry `zeed "All changelog entries are saved within <your_project_dir>/.zeed/"`
- Unify the entries `zeed unify`
- Unify then delete the entries `zeed unify --flush`
- Copy/Paste the unified entries in your current changelog file

#### How to work with channels?
Each entry is related to a channel. The default channel is `undefined`. To add support for a channel to your project, edit `.zeed/.zeed.yaml` file that way:
``` yaml
channels:
  - added
  - changed
  - deprecated
  - fixed
  - removed
  - security
  - name_of_your_channel # Only a-z and _ are allowed
```
Channels serve to group your entries and it will be useful for templates.

#### How to work with templates?
Templates serve to customize the rendering of `unify` command.

Templates must use [Go templating engine](https://golang.org/pkg/text/template/). They have access to this fields :
- Entries (list of all entries)
- Channels (list of all entries grouped by channel)

Add your own template by following this steps:
- create a file within `.zeed` directory
  - Example of content:
  ```
  {{range .Entries -}}
  {{- if eq .Channel.Id "undefined" -}}
  - [Undefined] {{.Text}} ({{.Priority}})
  {{- end}}
  {{- if eq .Channel.Id "added" -}}
  - [Added] {{.Text}} ({{.Priority}})
  {{- end}}
  {{end -}}
  ```
- Use your new file as a template when you unify `zeed unify --template filename`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We will use [SemVer](http://semver.org/) for versioning.

## Maintainers

* **Souhail** - [Profile](https://github.com/souhail-5/)

See also the list of [contributors](https://github.com/souhail-5/zeed/graphs/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.