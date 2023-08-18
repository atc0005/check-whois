// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/atc0005/check-whois/internal/domain"
	"github.com/atc0005/go-nagios"
)

// getPerfData generates performance data metrics from the given domain
// metadata age thresholds. An error is returned if any are encountered while
// gathering metrics or if invalid domain metadata is provided.
func getPerfData(d *domain.Metadata, ageCritical int, ageWarning int) ([]nagios.PerformanceData, error) {

	var expires int
	if daysToExpiration, err := domain.UntilExpiration(d); err == nil {
		expires = daysToExpiration
	}

	var updated int
	if daysSinceUpdate, err := domain.SinceUpdate(d); err == nil {
		updated = daysSinceUpdate
	}

	var created int
	if daysSinceCreation, err := domain.SinceCreation(d); err == nil {
		created = daysSinceCreation
	}

	pd := []nagios.PerformanceData{
		{
			Label:             "expires",
			Value:             fmt.Sprintf("%d", expires),
			UnitOfMeasurement: "d",
			Warn:              fmt.Sprintf("%d", ageWarning),
			Crit:              fmt.Sprintf("%d", ageCritical),
		},
		{
			Label:             "since_update",
			Value:             fmt.Sprintf("%d", updated),
			UnitOfMeasurement: "d",
		},
		{
			Label:             "since_creation",
			Value:             fmt.Sprintf("%d", created),
			UnitOfMeasurement: "d",
		},
	}

	return pd, nil

}
