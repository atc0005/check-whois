// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package whois

import (
	"fmt"
	"strings"
	"time"

	"github.com/atc0005/check-whois/internal/domain/metadata"

	whoisparser "github.com/likexian/whois-parser"
)

// New instantiates a new Metadata type from parsed WHOIS data.
func New(whoisInfo *whoisparser.WhoisInfo, ageWarning time.Time, ageCritical time.Time) (*metadata.Domain, error) {
	var expirationDate time.Time
	var updatedDate time.Time
	var createdDate time.Time

	var err error

	// We attempt to use an already parsed time value as-is first, but if not
	// set we perform a cursory parsing attempt against the plaintext version
	// of the date values recorded in the parsed WHOIS record.
	switch {
	case whoisInfo.Domain.ExpirationDateInTime != nil:
		expirationDate = *whoisInfo.Domain.ExpirationDateInTime
	default:
		expirationDate, err = metadata.ParseDateString(whoisInfo.Domain.ExpirationDate)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse domain expiration date: %w",
				err,
			)
		}
	}

	switch {
	case whoisInfo.Domain.UpdatedDateInTime != nil:
		updatedDate = *whoisInfo.Domain.UpdatedDateInTime
	default:
		updatedDate, err = metadata.ParseDateString(whoisInfo.Domain.UpdatedDate)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse domain updated date: %w",
				err,
			)
		}
	}

	switch {
	case whoisInfo.Domain.CreatedDateInTime != nil:
		createdDate = *whoisInfo.Domain.CreatedDateInTime
	default:
		createdDate, err = metadata.ParseDateString(whoisInfo.Domain.CreatedDate)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse domain creation date: %w",
				err,
			)
		}
	}

	d := metadata.Domain{
		AgeWarningThreshold:  ageWarning,
		AgeCriticalThreshold: ageCritical,
		Name:                 domainName(whoisInfo),
		Status:               domainStatus(whoisInfo),
		DataSource:           metadata.DataSourceWHOIS,
		RegistrarName:        registrarName(whoisInfo),
		RegistrantName:       registrantName(whoisInfo),
		RegistrantEmail:      registrantEmail(whoisInfo),
		ExpirationDate:       expirationDate,
		UpdatedDate:          updatedDate,
		CreatedDate:          createdDate,
	}

	return &d, nil

}

// domainName provides the domain name value from the WhoIS record or the
// fallback/placeholder value for the field.
func domainName(w *whoisparser.WhoisInfo) string {
	if w.Domain != nil {
		return w.Domain.Domain
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// domainStatus provides the domain status value from the WhoIS record or the
// fallback/placeholder value for the field.
func domainStatus(w *whoisparser.WhoisInfo) string {
	if w.Domain != nil && len(w.Domain.Status) != 0 {
		return strings.Join(w.Domain.Status, ", ")
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrarName provides the registrar name value from the WhoIS record or
// the fallback/placeholder value for the field.
func registrarName(w *whoisparser.WhoisInfo) string {
	if w.Registrar != nil && w.Registrar.Name != "" {
		return w.Registrar.Name
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrantName provides the registrant name value from the WhoIS record or
// the fallback/placeholder value for the field.
func registrantName(w *whoisparser.WhoisInfo) string {
	if w.Registrant != nil && w.Registrant.Name != "" {
		return w.Registrant.Name
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrantEmail provides the registrant email value from the WhoIS record
// or the fallback/placeholder value for the field.
func registrantEmail(w *whoisparser.WhoisInfo) string {
	if w.Registrant != nil && w.Registrant.Email != "" {
		return w.Registrant.Email
	}

	return metadata.DefaultWhoISPlaceholderValue
}
