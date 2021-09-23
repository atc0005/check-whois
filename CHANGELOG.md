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

[Unreleased]: https://github.com/atc0005/check-whois/compare/v0.2.2...HEAD
[v0.2.2]: https://github.com/atc0005/check-whois/releases/tag/v0.2.2
[v0.2.1]: https://github.com/atc0005/check-whois/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/check-whois/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/check-whois/releases/tag/v0.1.0
