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

	// Set initial "state" as valid, adjust as we go.
	var nagiosExitState = nagios.ExitState{
		LastError:      nil,
		ExitStatusCode: nagios.StateOKExitCode,
	}

	// defer this from the start so it is the last deferred function to run
	defer nagiosExitState.ReturnCheckResults()

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
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error initializing application",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.LastError = cfgErr
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	// Use provided threshold values to calculate the expiration times that
	// should trigger either a WARNING or CRITICAL state.
	now := time.Now().UTC()
	domainExpireAgeWarning := now.AddDate(0, 0, cfg.AgeWarning)
	domainExpireAgeCritical := now.AddDate(0, 0, cfg.AgeCritical)

	nagiosExitState.WarningThreshold = fmt.Sprintf(
		"Expires before %v (%d days)",
		domainExpireAgeWarning.Format(domain.DomainDateLayout),
		cfg.AgeWarning,
	)
	nagiosExitState.CriticalThreshold = fmt.Sprintf(
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

		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error fetching WHOIS data for %s domain",
			nagios.StateCRITICALLabel,
			cfg.Domain,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	}

	parsedWhois, err := whoisparser.Parse(whoisRaw)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse WHOIS data")

		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error parsing WHOIS data for %s domain",
			nagios.StateCRITICALLabel,
			cfg.Domain,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	}

	d, err := domain.NewDomain(parsedWhois, domainExpireAgeWarning, domainExpireAgeCritical)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse WhoisInfo data")

		nagiosExitState.LastError = err
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error parsing WhoisInfo data for %s domain",
			nagios.StateCRITICALLabel,
			cfg.Domain,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	}

	switch {

	case d.IsExpired():

		log.Error().Msg("Domain has expired")

		nagiosExitState.LastError = domain.ErrDomainExpired
		nagiosExitState.ServiceOutput = d.OneLineCheckSummary()
		nagiosExitState.LongServiceOutput = d.Report()
		nagiosExitState.ExitStatusCode = d.ServiceState().ExitCode

		return

	case d.IsExpiring():

		log.Warn().Msg("Domain is expiring")

		nagiosExitState.LastError = domain.ErrDomainExpiring
		nagiosExitState.ServiceOutput = d.OneLineCheckSummary()
		nagiosExitState.LongServiceOutput = d.Report()
		nagiosExitState.ExitStatusCode = d.ServiceState().ExitCode

		return

	default:

		log.Debug().Msg("No problems with expiration date for domain detected")

		nagiosExitState.LastError = nil
		nagiosExitState.ServiceOutput = d.OneLineCheckSummary()
		nagiosExitState.LongServiceOutput = d.Report()
		nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

		return

	}

}
