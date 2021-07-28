// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

const myAppName string = "check-whois"
const myAppURL string = "https://github.com/atc0005/" + myAppName

const (
	domainFlagHelp                  string = "The name of the domain whose WHOIS records will be evaluated."
	versionFlagHelp                 string = "Whether to display application version and then immediately exit application."
	logLevelFlagHelp                string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	brandingFlagHelp                string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	domainExpireAgeWarningFlagHelp  string = "The number of days remaining before domain expiration when a WARNING state is triggered."
	domainExpireAgeCriticalFlagHelp string = "The number of days remaining before domain expiration when a CRITICAL state is triggered."
)

// Default flag settings if not overridden by user input
const (
	defaultDomain                string = ""
	defaultLogLevel              string = "info"
	defaultBranding              bool   = false
	defaultDisplayVersionAndExit bool   = false

	// Default WARNING threshold is 30 days
	defaultDomainExpireAgeWarning int = 30

	// Default CRITICAL threshold is 15 days
	defaultDomainExpireAgeCritical int = 15
)

const (

	// LogLevelDisabled maps to zerolog.Disabled logging level
	LogLevelDisabled string = "disabled"

	// LogLevelPanic maps to zerolog.PanicLevel logging level
	LogLevelPanic string = "panic"

	// LogLevelFatal maps to zerolog.FatalLevel logging level
	LogLevelFatal string = "fatal"

	// LogLevelError maps to zerolog.ErrorLevel logging level
	LogLevelError string = "error"

	// LogLevelWarn maps to zerolog.WarnLevel logging level
	LogLevelWarn string = "warn"

	// LogLevelInfo maps to zerolog.InfoLevel logging level
	LogLevelInfo string = "info"

	// LogLevelDebug maps to zerolog.DebugLevel logging level
	LogLevelDebug string = "debug"

	// LogLevelTrace maps to zerolog.TraceLevel logging level
	LogLevelTrace string = "trace"
)
