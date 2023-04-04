// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"errors"
	"fmt"
	"time"

	zlog "github.com/rs/zerolog/log"

	"github.com/atc0005/check-whois/internal/config"
	"github.com/atc0005/check-whois/internal/domain"

	"github.com/atc0005/go-nagios"
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func main() {

	plugin := nagios.NewPlugin()

	// defer this from the start so it is the last deferred function to run
	defer plugin.ReturnCheckResults()

	// Setup configuration by parsing user-provided flags.
	cfg, cfgErr := config.New()
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case cfgErr != nil:
		// We're using the standalone Err function from rs/zerolog/log as we
		// do not have a working configuration.
		zlog.Err(cfgErr).Msg("Error initializing application")
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error initializing application",
			nagios.StateUNKNOWNLabel,
		)
		plugin.AddError(cfgErr)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode

		return
	}

	// Use provided threshold values to calculate the expiration times that
	// should trigger either a WARNING or CRITICAL state.
	now := time.Now().UTC()
	domainExpireAgeWarning := now.AddDate(0, 0, cfg.AgeWarning)
	domainExpireAgeCritical := now.AddDate(0, 0, cfg.AgeCritical)

	plugin.WarningThreshold = fmt.Sprintf(
		"Expires before %v (%d days)",
		domainExpireAgeWarning.Format(domain.DomainDateLayout),
		cfg.AgeWarning,
	)
	plugin.CriticalThreshold = fmt.Sprintf(
		"Expires before %v (%d days)",
		domainExpireAgeCritical.Format(domain.DomainDateLayout),
		cfg.AgeCritical,
	)

	log := cfg.Log.With().
		Str("domain", cfg.Domain).
		Logger()

	var whoisRaw string
	var err error
	switch {
	case cfg.RegistrarServer != "":
		whoisRaw, err = whois.Whois(cfg.Domain, cfg.RegistrarServer)
	default:
		whoisRaw, err = whois.Whois(cfg.Domain)
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to query WHOIS data")

		plugin.AddError(err)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error fetching WHOIS data for %s domain",
			nagios.StateUNKNOWNLabel,
			cfg.Domain,
		)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode

		return

	}

	parsedWhois, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse WHOIS data")

		plugin.AddError(err)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error parsing WHOIS data for %s domain",
			nagios.StateUNKNOWNLabel,
			cfg.Domain,
		)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode

		return

	}

	d, err := domain.NewDomain(parsedWhois, domainExpireAgeWarning, domainExpireAgeCritical)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse WhoisInfo data")

		plugin.AddError(err)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error parsing WhoisInfo data for %s domain",
			nagios.StateUNKNOWNLabel,
			cfg.Domain,
		)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode

		return

	}

	switch {

	case d.IsExpired():

		log.Error().Msg("Domain has expired")

		plugin.AddError(domain.ErrDomainExpired)
		plugin.ServiceOutput = d.OneLineCheckSummary()
		plugin.LongServiceOutput = d.Report()
		plugin.ExitStatusCode = d.ServiceState().ExitCode

		return

	case d.IsExpiring():

		log.Warn().Msg("Domain is expiring")

		plugin.AddError(domain.ErrDomainExpiring)
		plugin.ServiceOutput = d.OneLineCheckSummary()
		plugin.LongServiceOutput = d.Report()
		plugin.ExitStatusCode = d.ServiceState().ExitCode

		return

	default:

		log.Debug().Msg("No problems with expiration date for domain detected")

		plugin.ServiceOutput = d.OneLineCheckSummary()
		plugin.LongServiceOutput = d.Report()
		plugin.ExitStatusCode = nagios.StateOKExitCode

		return

	}

}
