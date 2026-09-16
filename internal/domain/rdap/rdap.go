// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package rdap

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/atc0005/check-whois/internal/domain/metadata"

	"github.com/openrdap/rdap"
)

// New instantiates a new Metadata type from parsed RDAP domain data.
func New(domain *rdap.Domain, bootstrapRegistryFile []byte, ageWarning time.Time, ageCritical time.Time) (*metadata.Domain, error) {
	var expirationDate time.Time
	var updatedDate time.Time
	var createdDate time.Time

	var err error

	var bootstrapMetadata metadata.RDAPDNSBootstrap
	if err := json.Unmarshal(bootstrapRegistryFile, &bootstrapMetadata); err != nil {
		return nil, fmt.Errorf(
			"unable to parse RDAP bootstrap file for Domain Name System registrations: %w",
			metadata.ErrMissingValue,
		)

	}

	if domain == nil {
		return nil, fmt.Errorf(
			"unable to evaluate domain metadata: %w",
			metadata.ErrMissingValue,
		)
	}

	for _, event := range domain.Events {
		if event.Action == "expiration" {
			expirationDate, err = time.Parse(time.RFC3339, event.Date)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse domain expiration date: %w",
					err,
				)
			}
		}
	}

	for _, event := range domain.Events {
		if event.Action == "last changed" {
			updatedDate, err = time.Parse(time.RFC3339, event.Date)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse domain updated date: %w",
					err,
				)
			}
		}
	}

	for _, event := range domain.Events {
		if event.Action == "registration" {
			createdDate, err = time.Parse(time.RFC3339, event.Date)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to parse domain creation date: %w",
					err,
				)
			}
		}
	}

	d := metadata.Domain{
		AgeWarningThreshold:  ageWarning,
		AgeCriticalThreshold: ageCritical,
		Name:                 domain.LDHName,
		Status:               domainStatus(domain),
		DataSource:           metadata.DataSourceRDAP,
		RDAPDNSBootstrap:     bootstrapMetadata,
		RegistrarName:        registrarName(domain),
		RegistrantName:       registrantName(domain),
		RegistrantEmail:      registrantEmail(domain),
		ExpirationDate:       expirationDate,
		UpdatedDate:          updatedDate,
		CreatedDate:          createdDate,
	}

	return &d, nil

}

// domainStatus provides the domain status value from the RDAP record or the
// fallback/placeholder value for the field.
func domainStatus(d *rdap.Domain) string {
	if len(d.Status) != 0 {
		// TODO: Might need to handle deduplication of status values.
		return strings.Join(d.Status, ", ")
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrarName provides the registrar name value from the RDAP record or
// the fallback/placeholder value for the field.
func registrarName(d *rdap.Domain) string {
	entity := findEntity("registrar", d.Entities)
	if entity == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	v := entity.VCard
	if v == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	if v.Name() != "" {
		return v.Name()
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrantName provides the registrant name value from the RDAP record or
// the fallback/placeholder value for the field.
func registrantName(d *rdap.Domain) string {
	entity := findEntity("registrant", d.Entities)
	if entity == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	v := entity.VCard
	if v == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	if v.Name() != "" {
		return v.Name()
	}

	return metadata.DefaultWhoISPlaceholderValue
}

// registrantEmail provides the registrant email value from the RDAP record
// or the fallback/placeholder value for the field.
func registrantEmail(d *rdap.Domain) string {
	entity := findEntity("registrant", d.Entities)
	if entity == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	v := entity.VCard
	if v == nil {
		return metadata.DefaultWhoISPlaceholderValue
	}

	if v.Email() != "" {
		return v.Email()
	}

	return metadata.DefaultWhoISPlaceholderValue
}

func findEntity(role string, entities []rdap.Entity) *rdap.Entity {
	for _, e := range entities {
		if slices.Contains(e.Roles, role) {
			return &e
		}
	}

	return nil
}
