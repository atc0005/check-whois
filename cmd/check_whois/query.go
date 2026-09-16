// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/atc0005/check-whois/internal/config"
	"github.com/atc0005/check-whois/internal/domain/metadata"
	rdapDomain "github.com/atc0005/check-whois/internal/domain/rdap"
	whoisDomain "github.com/atc0005/check-whois/internal/domain/whois"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/openrdap/rdap"
	"github.com/openrdap/rdap/bootstrap"
	"github.com/rs/zerolog"
)

// Sourced periodically from https://data.iana.org/rdap/dns.json
//
// See also:
// - https://deployment.rdap.org/
// - https://www.iana.org/assignments/rdap-dns
// - https://about.rdap.org/#additional
//
//go:embed dns.json
var bootstrapRegistryFile []byte

// mustUseWHOIS indicates whether a WHOIS query (deprecated) is necessary to
// obtain information for a given domain or whether a RDAP query (default)
// will be used.
func mustUseWHOIS(domain string, rdapClient *rdap.Client) bool {
	question := &bootstrap.Question{
		RegistryType: bootstrap.DNS,
		Query:        strings.ToLower(domain),
	}

	answer, err := rdapClient.Bootstrap.Lookup(question)
	if err != nil || len(answer.URLs) == 0 {
		// Failed to lookup RDAP base URLs, so fallback to WHOIS query.
		return true
	}

	return false
}

// RDAPDebugLogger provides a wrapper for enabling verbose logging from the
// openrdap/rdap library.
func RDAPDebugLogger(log *zerolog.Logger) func(msg string) {
	return func(msg string) {
		log.Debug().Msgf("openrdap/rdap: %s", msg)
	}
}

func rdapQuery(cfg *config.Config, rdapClient *rdap.Client, ageWarning time.Time, ageCritical time.Time) (*metadata.Domain, error) {
	// Define the domain query.
	req := &rdap.Request{
		Type:  rdap.DomainRequest,
		Query: cfg.Domain,
	}

	if cfg.RDAPServerURL != "" {
		u, err := url.Parse(cfg.RDAPServerURL)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse custom RDAP Server URL: %w", err,
			)
		}

		req = req.WithServer(u)
	}

	// Execute the RDAP query.
	resp, err := rdapClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request RDAP data: %w", err)
	}

	// Cast the response object to a Domain type.
	rd, ok := resp.Object.(*rdap.Domain)
	if !ok {
		return nil, errors.New(
			"failed to convert response object to domain type",
		)
	}

	domain, err := rdapDomain.New(rd, bootstrapRegistryFile, ageWarning, ageCritical)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain metadata: %w", err,
		)
	}

	return domain, nil
}

func whoisQuery(cfg *config.Config, ageWarning time.Time, ageCritical time.Time) (*metadata.Domain, error) {
	var whoisRaw string
	var err error

	client := whois.NewClient()

	// Explicitly set referral lookup behavior. Referral lookups are performed
	// unless requested otherwise by the sysadmin.
	client.SetDisableReferral(cfg.DisableReferralLookups)

	switch {
	case cfg.RegistrarServer != "":
		whoisRaw, err = client.Whois(cfg.Domain, cfg.RegistrarServer)
	default:
		whoisRaw, err = client.Whois(cfg.Domain)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query WHOIS data: %w", err)
	}

	parsedWhois, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse WHOIS data: %w", err)
	}

	domain, err := whoisDomain.New(&parsedWhois, ageWarning, ageCritical)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain metadata: %w", err,
		)
	}

	return domain, nil
}
