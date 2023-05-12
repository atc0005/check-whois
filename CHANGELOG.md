# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/check-whois/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.3.1] - 2023-05-12

### Overview

- Dependency updates
- built using Go 1.19.9
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.7` to `1.19.9`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.3` to `go-ci-oldstable-build-v0.10.5`
  - `github.com/likexian/whois`
    - `v1.14.6` to `v1.50.0`
  - `github.com/likexian/whois-parser`
    - `v1.24.7` to `v1.24.8`
  - `github.com/likexian/gokit`
    - `v0.25.11` to `v0.25.13`
  - `rs/zerolog`
    - `v1.29.0` to `v1.29.1`
  - `golang.org/x/net`
    - `v0.8.0` to `v0.10.0`
  - `golang.org/x/sys`
    - `v0.6.0` to `v0.8.0`
  - `golang.org/x/text`
    - `v0.8.0` to `v0.9.0`

## [v0.3.0] - 2023-04-04

### Overview

- Add support for generating DEB, RPM packages
- Build improvements
- Generated binary changes
  - filename patterns
  - compression (~ 66% smaller)
  - executable metadata
- built using Go 1.19.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-154) Generate RPM/DEB packages using nFPM
- (GH-155) Add version details to Windows executables

### Changed

- (GH-153) Switch to semantic versioning (semver) compatible versioning
  pattern
- (GH-156) Makefile: Compress binaries & use fixed filenames
- (GH-157) Makefile: Refresh recipes to add "standard" set, new
  package-related options
- (GH-158) Build dev/stable releases using go-ci Docker image

## [v0.2.14] - 2023-04-04

### Overview

- Bug fixes
- GitHub Actions workflows updates
- Dependency updates
- built using Go 1.19.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- (GH-136) Add Go Module Validation, Dependency Updates jobs

### Changed

- Dependencies
  - `Go`
    - `1.19.4` to `1.19.7`
  - `github.com/likexian/whois`
    - `v1.14.4` to `v1.14.6`
  - `github.com/likexian/whois-parser`
    - `v1.24.2` to `v1.24.7`
  - `github.com/likexian/gokit`
    - `v0.25.9` to `v0.25.11`
  - `atc0005/go-nagios`
    - `v0.10.2` to `v0.14.0`
  - `rs/zerolog`
    - `v1.28.0` to `v1.29.0`
  - `github.com/mattn/go-isatty`
    - `v0.0.16` to `v0.0.18`
  - `golang.org/x/net`
    - `v0.4.0` to `v0.8.0`
  - `golang.org/x/sys`
    - `v0.3.0` to `v0.6.0`
  - `golang.org/x/text`
    - `v0.5.0` to `v0.8.0`
- Misc
  - (GH-130) Update nagios library usage
- CI
  - (GH-145) Drop `Push Validation` workflow
  - (GH-146) Rework workflow scheduling
  - (GH-148) Remove `Push Validation` workflow status badge

### Fixed

- (GH-167) Update vuln analysis GHAW to use on.push
- (GH-170) Use UNKNOWN state for invalid command-line options
- (GH-171) Use UNKNOWN state for fetching & parsing failures

## [v0.2.13] - 2022-12-09

### Overview

- Dependency updates
- built using Go 1.19.4
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.3` to `1.19.4`
  - `github.com/likexian/whois`
    - `v1.14.3-0.20221012013111-a48608e6097a` to `v1.14.4`
  - `github.com/likexian/whois-parser`
    - `v1.24.1` to `v1.24.2`
  - `github.com/mattn/go-colorable`
    - `v0.1.12` to `v0.1.13`
  - `github.com/mattn/go-isatty`
    - `v0.0.14` to `v0.0.16`
  - `golang.org/x/net`
    - `v0.0.0-20220708220712-1185a9018129` to `v0.4.0`
  - `golang.org/x/sys`
    - `v0.0.0-20220520151302-bc2c85ada10a` to `v0.3.0`
  - `golang.org/x/text`
    - `v0.3.7` to `v0.5.0`

## [v0.2.12] - 2022-11-09

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.3
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.19.1` to `1.19.3`
  - `github.com/likexian/whois`
    - `v1.14.2` to `v1.14.3-0.20221012013111-a48608e6097a`
  - `atc0005/go-nagios`
    - `v0.10.0` to `v0.10.2`
- (GH-108) Refactor GitHub Actions workflows to import logic

### Fixed

- (GH-115) Receiving `dial tcp: address tcp///whois.pairdomains.com: unknown
  port` error for a specific domain
- (GH-117) WHOIS CreatedDate value incorrectly duplicates ExpirationDate value
- (GH-119) Failure to parse date fields in WHOIS record: `parsing time
  "2022-06-03" as "2006-01-02T15:04:05Z07:00": cannot parse "" as "T"`
- (GH-121) Fix Makefile Go module base path detection

## [v0.2.11] - 2022-09-22

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.1
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.13` to `1.19.1`
  - `atc0005/go-nagios`
    - `v0.9.1` to `v0.10.0`
  - `rs/zerolog`
    - `v1.27.0` to `v1.28.0`
  - `github/codeql-action`
    - `v2.1.21` to `v2.1.25`
- (GH-98) Update project to Go 1.19
- (GH-99) Update Makefile and GitHub Actions Workflows

## [v0.2.10] - 2022-08-25

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.17.13
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.12` to `1.17.13`

### Fixed

- (GH-95) Apply Go 1.19 specific doc comments linting fixes

## [v0.2.9] - 2022-07-21

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.17.12
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.10` to `1.17.12`
  - `likexian/whois`
    - `v1.12.5` to `v1.14.2`
  - `likexian/whois-parser`
    - `v1.22.0` to `v1.24.1`
  - `atc0005/go-nagios`
    - `v0.8.2` to `v0.9.1`
  - `rs/zerolog`
    - `v1.26.1` to `v1.27.0`

### Fixed

- (GH-88) Update lintinstall Makefile recipe
- (GH-90) Fix atc0005/go-nagios library usage linting errors
- (GH-93) Fix unused Markdownlint reference

## [v0.2.8] - 2022-05-13

### Overview

- Dependency updates
- built using Go 1.17.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.9` to `1.17.10`

## [v0.2.7] - 2022-04-29

### Overview

- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.7` to `1.17.9`
  - `likexian/whois`
    - `v1.12.4` to `v1.12.5`

## [v0.2.6] - 2022-03-02

### Overview

- Dependency updates
- CI / linting improvements
- built using Go 1.17.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.17.6` to `1.17.7`
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

### Fixed

- (GH-66) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-67) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)

## [v0.2.5] - 2022-01-21

### Overview

- Dependency updates
- built using Go 1.17.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.12` to `1.17.6`
    - (GH-60) Update go.mod file, canary Dockerfile to reflect current
      dependencies
  - `atc0005/go-nagios`
    - `v0.8.1` to `v0.8.2`

### Fixed

- (GH-61) Remove additional cert references
- (GH-62) Fix doc comments for FormattedExpiration func

## [v0.2.4] - 2021-12-28

### Overview

- Dependency updates
- built using Go 1.16.12
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.10` to `1.16.12`
  - `rs/zerolog`
    - `v1.26.0` to `v1.26.1`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.1`

- (GH-54) Help output generated by `-h`, `--help` flag is sent to `stderr`,
  should go to `stdout` instead

### Fixed

- (GH-49) Add missing subsection to CHANGELOG entry

## [v0.2.3] - 2021-11-08

### Overview

- Dependency updates
- built using Go 1.16.10
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.8` to `1.16.10`
  - `atc0005/go-nagios`
    - `v0.7.0` to `v0.8.1`
  - `rs/zerolog`
    - `v1.25.0` to `v1.26.0`
  - `actions/checkout`
    - `v2.3.4` to `v2.4.0`
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`
  - `likexian/whois`
    - `v1.12.1` to `v1.12.4`
  - `likexian/whois-parser`
    - `v1.20.5` to `v1.22.0`

- (GH-39) Replace bundled ServiceState type

### Fixed

- (GH-41) Fix various typos and copy/paste/modify issues

## [v0.2.2] - 2021-09-23

### Overview

- Dependency updates
- built using Go 1.16.8
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.8`
  - `atc0005/go-nagios`
    - `v0.6.1` to `v0.7.0`
  - `rs/zerolog`
    - `v1.23.0` to `v1.25.0`
  - `likexian/whois-parser`
    - `v1.20.4` to `v1.20.5`

## [v0.2.1] - 2021-08-06

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - updated from `v2.3.0` to `v2.4.0`

## [v0.2.0] - 2021-07-30

### Overview

- Add new optional flag
- built using Go 1.16.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- Add optional flag for registrar WHOIS server

## [v0.1.0] - 2021-07-30

### Overview

- Initial project release
- built using Go 1.16.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

This release provides an early version of a Nagios plugin used to check the
expiration date for a specified domain. See the README file for additional
details.

Features of this release:

- Nagios plugin for monitoring expiration of WHOIS records

- Optional branding "signature"
  - used to indicate what Nagios plugin (and what version) is responsible for
    the service check result

- Optional, leveled logging using `rs/zerolog` package
  - JSON-format output (to `stderr`)
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`.

[Unreleased]: https://github.com/atc0005/check-whois/compare/v0.3.1...HEAD
[v0.3.1]: https://github.com/atc0005/check-whois/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/check-whois/releases/tag/v0.3.0
[v0.2.14]: https://github.com/atc0005/check-whois/releases/tag/v0.2.14
[v0.2.13]: https://github.com/atc0005/check-whois/releases/tag/v0.2.13
[v0.2.12]: https://github.com/atc0005/check-whois/releases/tag/v0.2.12
[v0.2.11]: https://github.com/atc0005/check-whois/releases/tag/v0.2.11
[v0.2.10]: https://github.com/atc0005/check-whois/releases/tag/v0.2.10
[v0.2.9]: https://github.com/atc0005/check-whois/releases/tag/v0.2.9
[v0.2.8]: https://github.com/atc0005/check-whois/releases/tag/v0.2.8
[v0.2.7]: https://github.com/atc0005/check-whois/releases/tag/v0.2.7
[v0.2.6]: https://github.com/atc0005/check-whois/releases/tag/v0.2.6
[v0.2.5]: https://github.com/atc0005/check-whois/releases/tag/v0.2.5
[v0.2.4]: https://github.com/atc0005/check-whois/releases/tag/v0.2.4
[v0.2.3]: https://github.com/atc0005/check-whois/releases/tag/v0.2.3
[v0.2.2]: https://github.com/atc0005/check-whois/releases/tag/v0.2.2
[v0.2.1]: https://github.com/atc0005/check-whois/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/check-whois/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/check-whois/releases/tag/v0.1.0
