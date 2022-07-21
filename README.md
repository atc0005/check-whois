<!-- omit in toc -->
# check-whois

Go-based tooling to monitor WHOIS records.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-whois.svg?style=flat-square)][repo-url]
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-whois.svg)](https://pkg.go.dev/github.com/atc0005/check-whois)
[![Validate Codebase](https://github.com/atc0005/check-whois/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/check-whois/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/check-whois/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/check-whois/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/check-whois/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/check-whois/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/check-whois/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/check-whois/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
  - [`check_whois`](#check_whois)
- [Features](#features)
- [Changelog](#changelog)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using release binaries](#using-release-binaries)
- [Configuration](#configuration)
  - [Command-line arguments](#command-line-arguments)
    - [`check_whois`](#check_whois-1)
- [Examples](#examples)
  - [`OK` result](#ok-result)
  - [`WARNING` result](#warning-result)
  - [`CRITICAL` result](#critical-result)
- [License](#license)
- [Related projects](#related-projects)
- [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

## Overview

This repo is intended to provide various tools used to monitor WHOIS.

| Tool Name     | Overall Status | Description                                               |
| ------------- | -------------- | --------------------------------------------------------- |
| `check_whois` | Alpha          | Nagios plugin used to monitor expiration of WHOIS records |

### `check_whois`

Nagios plugin used to monitor expiration of WHOIS records.

The output for this application is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for use in email and Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

## Features

- Nagios plugin for monitoring expiration of WHOIS records

- Optional branding "signature"
  - used to indicate what Nagios plugin (and what version) is responsible for
    the service check result

- Optional, leveled logging using `rs/zerolog` package
  - JSON-format output (to `stderr`)
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`.

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go
  - see this project's `go.mod` file for *preferred* version
  - this project tests against [officially supported Go
    releases][go-supported-releases]
    - the most recent stable release (aka, "stable")
    - the prior, but still supported release (aka, "oldstable")
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 10
- Ubuntu Linux 18.04+

## Installation

### From source

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/check-whois`
   1. `cd check-whois`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     1. `sudo yum install make gcc`
1. Build
   - for current operating system
     - `go build -mod=vendor ./cmd/check_whois/`
       - *forces build to use bundled dependencies in top-level `vendor`
         folder*
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for Windows
      - `make windows`
   - for Linux
     - `make linux`
1. Locate generated binaries
   - if using `Makefile`
     - look in `/tmp/check-whois/release_assets/check_whois/`
   - if using `go build`
     - look in `/tmp/check-whois/`
1. Copy the applicable binaries to whatever systems needs to run them
1. Deploy
   - Place `list-emails` in a location of your choice
   - Place `check_whois` in the same location where your distro's
     package manage has place other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_whois` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_whois` on RedHat-based
       systems

### Using release binaries

1. Download the [latest
   release][repo-url] binaries
1. Deploy
   - Place `check_whois` in the same location where your distro's
     package manager places other Nagios plugins
     - as `/usr/lib/nagios/plugins/check_whois` on Debian-based systems
     - as `/usr/lib64/nagios/plugins/check_whois` on RedHat-based
       systems

## Configuration

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

#### `check_whois`

| Flag                | Required | Default | Repeat | Possible                                                                | Description                                                                                          |
| ------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `branding`          | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default. |
| `h`, `help`         | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                               |
| `v`, `version`      | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                        |
| `c`, `age-critical` | No       | 15      | No     | *positive whole number of days*                                         | The number of days remaining before domain expiration when a `CRITICAL` state is triggered.          |
| `w`, `age-warning`  | No       | 30      | No     | *positive whole number of days*                                         | The number of days remaining before domain expiration when a `WARNING` state is triggered.           |
| `ll`, `log-level`   | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                            |
| `d`, `domain`       | **Yes**  |         | No     | *domain name*                                                           | The name of the domain whose WHOIS records will be evaluated.                                        |
| `s`, `server`       | No       |         | No     | *valid WHOIS server fqdn*                                               | The name of the optional domain registrar WHOIS server to use for queries.                           |

## Examples

### `OK` result

This example uses default age thresholds to check the expiration date for the
specified domain. Since the expiration occurs after those default thresholds,
the result is considered `OK`.

```ShellSession
$ ./check_whois --domain "google.com"
OK: "google.com" domain registration has 2602d 18h remaining


**ERRORS**

* None

**THRESHOLDS**

* CRITICAL: Expires before 2021-08-14 09:31:39 +0000 UTC (15 days)
* WARNING: Expires before 2021-08-29 09:31:39 +0000 UTC (30 days)

**DETAILED INFO**

WHOIS metadata for "google.com" domain:

* Status: [clientdeleteprohibited, clienttransferprohibited, clientupdateprohibited, serverdeleteprohibited, servertransferprohibited, serverupdateprohibited]
* Creation Date: 2028-09-14 04:00:00 +0000 UTC
* Updated Date: 2019-09-09 15:39:04 +0000 UTC
* Expiration Date: 2028-09-14 04:00:00 +0000 UTC
* Registrar Name: MarkMonitor Inc.
* Registrant Name:
* Registrant Email: select request email form at https://domains.markmonitor.com/whois/google.com
```

### `WARNING` result

This example uses explicit age thresholds to simulate a domain that is
expiring soon, but hasn't yet crossed the final `CRITICAL` age threshold. The
result is considered to be a `WARNING` state.

```ShellSession
$ ./check_whois --domain "microsoft.com" --age-warning 365 --age-critical 120
{"level":"warn","version":"check-whois x.y.z (https://github.com/atc0005/check-whois)","logging_level":"info","domain":"microsoft.com","time":"2021-07-30T04:40:08-05:00","caller":"/mnt/t/github/check-whois/cmd/check_whois/main.go:142","message":"Domain is expiring"}
WARNING: "microsoft.com" domain registration has 276d 18h remaining


**ERRORS**

* domain is expiring

**THRESHOLDS**

* CRITICAL: Expires before 2021-11-27 09:40:08 +0000 UTC (120 days)
* WARNING: Expires before 2022-07-30 09:40:08 +0000 UTC (365 days)

**DETAILED INFO**

WHOIS metadata for "microsoft.com" domain:

* Status: [clientdeleteprohibited, clienttransferprohibited, clientupdateprohibited, serverdeleteprohibited, servertransferprohibited, serverupdateprohibited]
* Creation Date: 2022-05-03 04:00:00 +0000 UTC
* Updated Date: 2021-03-12 23:25:32 +0000 UTC
* Expiration Date: 2022-05-03 04:00:00 +0000 UTC
* Registrar Name: MarkMonitor Inc.
* Registrant Name: Domain Administrator
* Registrant Email: admin@domains.microsoft
```

### `CRITICAL` result

This example uses explicit age thresholds to simulate a domain that is
expiring soon and has crossed the final `CRITICAL` age threshold. The
result is considered to be a `CRITICAL` state.

The domain is expected to expire soon without direct intervention from the
current domain owner.

```ShellSession
$ ./check_whois --domain "godaddy.com" --age-warning 365 --age-critical 120
{"level":"warn","version":"check-whois x.y.z (https://github.com/atc0005/check-whois)","logging_level":"info","domain":"godaddy.com","time":"2021-07-30T04:41:06-05:00","caller":"/mnt/t/github/check-whois/cmd/check_whois/main.go:142","message":"Domain is expiring"}
CRITICAL: "godaddy.com" domain registration has 94d 2h remaining


**ERRORS**

* domain is expiring

**THRESHOLDS**

* CRITICAL: Expires before 2021-11-27 09:41:06 +0000 UTC (120 days)
* WARNING: Expires before 2022-07-30 09:41:06 +0000 UTC (365 days)

**DETAILED INFO**

WHOIS metadata for "godaddy.com" domain:

* Status: [clientdeleteprohibited, clientrenewprohibited, clienttransferprohibited, clientupdateprohibited, serverdeleteprohibited, servertransferprohibited, serverupdateprohibited]
* Creation Date: 2021-11-01 11:59:59 +0000 UTC
* Updated Date: 2020-04-07 14:26:27 +0000 UTC
* Expiration Date: 2021-11-01 11:59:59 +0000 UTC
* Registrar Name: GoDaddy.com, LLC
* Registrant Name:
* Registrant Email: select contact domain holder link at https://www.godaddy.com/whois/results.aspx?domain=godaddy.com
```

## License

```license
MIT License

Copyright 2021 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Related projects

- <https://github.com/atc0005/send2teams>
- <https://github.com/atc0005/check-whois>
- <https://github.com/atc0005/check-illiad>
- <https://github.com/atc0005/check-mail>
- <https://github.com/atc0005/check-path>
- <https://github.com/atc0005/check-vmware>
- <https://github.com/atc0005/nagios-debug>
- <https://github.com/atc0005/go-nagios>

## References

- <https://github.com/likexian/whois>
- <https://github.com/likexian/whois-parser>
- <https://github.com/rs/zerolog>
- <https://github.com/atc0005/go-nagios>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-whois>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[go-supported-releases]: <https://go.dev/doc/devel/release#policy> "Go Release Policy"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
