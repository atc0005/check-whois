// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

// handleFlagsConfig wraps flag setup code into a bundle for potential ease of
// use and future testability
func (c *Config) handleFlagsConfig() {

	flag.BoolVar(&c.EmitBranding, "branding", defaultBranding, brandingFlagHelp)

	flag.IntVar(&c.AgeWarning, "w", defaultDomainExpireAgeWarning, domainExpireAgeWarningFlagHelp)
	flag.IntVar(&c.AgeWarning, "age-warning", defaultDomainExpireAgeWarning, domainExpireAgeWarningFlagHelp)

	flag.IntVar(&c.AgeCritical, "c", defaultDomainExpireAgeCritical, domainExpireAgeCriticalFlagHelp)
	flag.IntVar(&c.AgeCritical, "age-critical", defaultDomainExpireAgeCritical, domainExpireAgeCriticalFlagHelp)

	flag.StringVar(&c.LoggingLevel, "ll", defaultLogLevel, logLevelFlagHelp)
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLogLevel, logLevelFlagHelp)

	flag.StringVar(&c.Domain, "d", defaultDomain, domainFlagHelp)
	flag.StringVar(&c.Domain, "domain", defaultDomain, domainFlagHelp)

	flag.BoolVar(&c.ShowVersion, "v", defaultDisplayVersionAndExit, versionFlagHelp)
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
