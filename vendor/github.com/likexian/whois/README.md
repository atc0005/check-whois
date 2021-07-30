# Whois

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/likexian/whois?status.svg)](https://godoc.org/github.com/likexian/whois)
[![Build Status](https://travis-ci.org/likexian/whois.svg?branch=master)](https://travis-ci.org/likexian/whois)
[![Go Report Card](https://goreportcard.com/badge/github.com/likexian/whois)](https://goreportcard.com/report/github.com/likexian/whois)
[![Code Cover](https://codecov.io/gh/likexian/whois/graph/badge.svg)](https://codecov.io/gh/likexian/whois)

Whois is a simple Go module for domain and ip whois information query.

## Overview

All of domain, IP include IPv4 and IPv6, ASN are supported.

You can directly using the binary distributions whois, follow [whois release tool](cmd/whois).

Or you can do development by using this golang module as below.

## Installation

```shell
go get -u github.com/likexian/whois
```

## Importing

```go
import (
    "github.com/likexian/whois"
)
```

## Documentation

Visit the docs on [GoDoc](https://godoc.org/github.com/likexian/whois)

## Example

### whois query for domain

```go
result, err := whois.Whois("likexian.com")
if err == nil {
    fmt.Println(result)
}
```

### whois query for IPv6

```go
result, err := whois.Whois("2001:dc7::1")
if err == nil {
    fmt.Println(result)
}
```

### whois query for IPv4

```go
result, err := whois.Whois("1.1.1.1")
if err == nil {
    fmt.Println(result)
}
```

### whois query for ASN

```go
// or whois.Whois("AS60614")
result, err := whois.Whois("60614")
if err == nil {
    fmt.Println(result)
}
```

## Whois information parsing

Please refer to [whois-parser](https://github.com/likexian/whois-parser)

## License

Copyright 2014-2021 [Li Kexian](https://www.likexian.com/)

Licensed under the Apache License 2.0

## Donation

If this project is helpful, please share it with friends.

If you want to thank me, you can [give me a cup of coffee](https://www.likexian.com/donate/).
