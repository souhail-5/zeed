# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/souhail-5/zeed/compare/master...develop)
- (nothing here)

## [1.0.0](https://github.com/souhail-5/zeed/compare/1.0.0-beta...1.0.0) - 2020-10-17
### Added
- Add a built-in template reproducing the format of keepachangelog.com
- Do not save entries with channels other than those configured
- Changelog filename should be configurable
- Add a verbose mode for each command

### Changed
- Rename &#34;undefined&#34; channel to &#34;default&#34;
- Rename entry&#39;s metadata: &#34;Priority&#34; to &#34;Weight&#34;
- Templates must be part of the config file
- Entry&#39;s filename should be only an ULID
- Replace nanoid by ulid

### Fixed
- Error strings should not be capitalized or end with punctuation
- Limit channel format to [a-z_]
- Remove all os.Exit()

## [1.0.0-beta](https://github.com/souhail-5/zeed/compare/1.0.0-alpha...1.0.0-beta) - 2020-10-27
### Added
- Add tests

### Changed
- Rename the tool/project: from "conflogt" to "zeed"

## 1.0.0-alpha - 2020-03
### Added
- Initial release
