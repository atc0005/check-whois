// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package metadata

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
)

// DomainDateLayout is the chosen date layout for displaying
// domain created/expiration date/time values across our application.
const DomainDateLayout string = "2006-01-02 15:04:05 -0700 MST"

// DefaultWhoISPlaceholderValue is used as a fallback value for any not found
// in the WhoIS record data (e.g., registrant name or email).
const DefaultWhoISPlaceholderValue string = "unspecified"

// Datasource constants to provide consistency throughout module tooling.
const (
	DataSourceWHOIS string = "WHOIS"
	DataSourceRDAP  string = "RDAP"
)

// ErrDomainExpired is returned whenever a specified domain has expired.
var ErrDomainExpired = errors.New("domain has expired")

// ErrDomainExpiring is returned whenever a specified domain is expiring.
var ErrDomainExpiring = errors.New("domain is expiring")

// ErrMissingValue indicates that an expected value was missing.
var ErrMissingValue = errors.New("missing expected value")

// RDAPDNSBootstrap provides the metadata (if applicable) for the RDAP
// bootstrap file for Domain Name System registrations.
type RDAPDNSBootstrap struct {
	Description string    `json:"description"`
	Publication time.Time `json:"publication"`
}

// Domain represents the details for a specified domain, including the
// expiration (age) thresholds used to determine plugin status and parsed
// values.
type Domain struct {
	// Name is the plaintext label for the domain.
	Name string

	// Status is the formatted collection of Domain Status values.
	Status string

	// DataSource indicates the source of retrieved domain metadata.
	DataSource string

	// RDAPDNSBootstrap provides the metadata (if applicable) for the RDAP
	// bootstrap file for Domain Name System registrations.
	RDAPDNSBootstrap RDAPDNSBootstrap

	// RegistrarName provides the registrar name value for the domain.
	RegistrarName string

	// RegistrantName provides the registrant name value for the domain.
	RegistrantName string

	// RegistrantEmail provides the registrant email value for the domain.
	RegistrantEmail string

	// ExpirationDate indicates when this domain expires.
	ExpirationDate time.Time

	// UpdatedDate indicates when the domain metadata was last updated.
	UpdatedDate time.Time

	// CreatedDate indicates when this domain was created/registered.
	CreatedDate time.Time

	// AgeWarningThreshold is the specified age threshold for when domains
	// with an expiration less than this value are considered to be in a
	// WARNING state.
	AgeWarningThreshold time.Time

	// AgeCRITICALThreshold is the specified age threshold for when domains
	// with an expiration less than this value are considered to be in a
	// CRITICAL state.
	AgeCriticalThreshold time.Time
}

// ParseDateString attempts to parse a given date string using detailed
// formats first, then falls back to a basic format before finally giving up
// and returning an error if no formats match.
func ParseDateString(dateString string) (time.Time, error) {
	var date time.Time
	var err error

	date, err = time.Parse(time.RFC3339, dateString)
	if err == nil {
		return date, nil
	}

	date, err = time.Parse("2006-01-02", dateString)
	if err == nil {
		return date, nil
	}

	// Use last encountered error as return value.
	return time.Time{}, fmt.Errorf(
		"failed to parse date string %s: %w",
		dateString,
		err,
	)
}

// OneLineCheckSummary generates a one-line summary of the domain WHOIS check
// results for display and notification purposes.
func (d Domain) OneLineCheckSummary() string {

	var summary string

	switch {
	case d.IsExpired():
		summary = fmt.Sprintf(
			"%s: %q domain registration EXPIRED %s%s",
			d.ServiceState().Label,
			d.Name,
			FormattedExpiration(d.ExpirationDate),
			nagios.CheckOutputEOL,
		)

	default:

		summary = fmt.Sprintf(
			"%s: %q domain registration has %s%s",
			d.ServiceState().Label,
			d.Name,
			FormattedExpiration(d.ExpirationDate),
			nagios.CheckOutputEOL,
		)

	}

	return summary

}

// Report provides an overview of domain details appropriate for display as
// the LongServiceOutput provided via the web UI or as email or Teams
// notifications.
func (d Domain) Report() string {

	var summary strings.Builder

	_, _ = fmt.Fprintf(
		&summary,
		"Metadata for %q domain:%s%s",
		d.Name,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Status: %s%s",
		d.Status,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Datasource: %s%s",
		d.DataSource,
		nagios.CheckOutputEOL,
	)

	if !d.RDAPDNSBootstrap.Publication.IsZero() {
		_, _ = fmt.Fprintf(
			&summary,
			"* RDAP bootstrap file (bundled) last updated: %s%s",
			d.RDAPDNSBootstrap.Publication.Format(time.RFC3339),
			nagios.CheckOutputEOL,
		)
	}

	_, _ = fmt.Fprintf(
		&summary,
		"* Creation Date: %v%s",
		d.CreatedDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Updated Date: %v%s",
		d.UpdatedDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Expiration Date: %v%s",
		d.ExpirationDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Registrar Name: %v%s",
		d.RegistrarName,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Registrant Name: %v%s",
		d.RegistrantName,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&summary,
		"* Registrant Email: %v%s",
		d.RegistrantEmail,
		nagios.CheckOutputEOL,
	)

	if d.DataSource == DataSourceRDAP && !d.RDAPDNSBootstrap.Publication.IsZero() {
		bestByDate := time.Now().Add(-90 * 24 * time.Hour)

		if d.RDAPDNSBootstrap.Publication.Before(bestByDate) {
			_, _ = fmt.Fprintf(
				&summary,
				"%s ⚠️ WARNING: The bundled RDAP bootstrap file last updated more than %s. RDAP search compatibility may be outdated.%s",
				nagios.CheckOutputEOL,
				humanizeDays(bestByDate, time.Now()),
				nagios.CheckOutputEOL,
			)
		}
	}

	return summary.String()
}

// IsExpired indicates whether the domain expiration date has passed.
func (d Domain) IsExpired() bool {
	return d.ExpirationDate.Before(time.Now())
}

// IsExpiring compares the domain's current expiration date against the
// provided CRITICAL and WARNING thresholds to determine if the domain is
// about to expire.
func (d Domain) IsExpiring() bool {

	switch {
	case !d.IsExpired() && d.ExpirationDate.Before(d.AgeCriticalThreshold):
		return true
	case !d.IsExpired() && d.ExpirationDate.Before(d.AgeWarningThreshold):
		return true
	}

	return false

}

// IsWarningState indicates whether a domain's expiration date has been
// determined to be in a WARNING state. This returns false if the expiration
// date is in an OK or CRITICAL state, true otherwise.
func (d Domain) IsWarningState() bool {
	if !d.IsExpired() &&
		d.ExpirationDate.Before(d.AgeWarningThreshold) &&
		!d.ExpirationDate.Before(d.AgeCriticalThreshold) {
		return true
	}

	return false

}

// IsCriticalState indicates whether a domain's expiration date has been
// determined to be in a CRITICAL state. This returns false if the expiration
// date is in an OK or WARNING state, true otherwise.
func (d Domain) IsCriticalState() bool {
	if d.IsExpired() || d.ExpirationDate.Before(d.AgeCriticalThreshold) {
		return true
	}

	return false

}

// IsOKState indicates whether a domain's expiration date has been determined
// to be in an OK state, without expired or expiring domain registration.
func (d Domain) IsOKState() bool {
	return !d.IsWarningState() && !d.IsCriticalState()
}

// ServiceState returns the appropriate Service Check Status label and exit
// code for the evaluated domain expiration metadata.
func (d Domain) ServiceState() nagios.ServiceState {

	var stateLabel string
	var stateExitCode int

	switch {
	case d.IsCriticalState():
		stateLabel = nagios.StateCRITICALLabel
		stateExitCode = nagios.StateCRITICALExitCode
	case d.IsWarningState():
		stateLabel = nagios.StateWARNINGLabel
		stateExitCode = nagios.StateWARNINGExitCode
	case d.IsOKState():
		stateLabel = nagios.StateOKLabel
		stateExitCode = nagios.StateOKExitCode
	default:
		stateLabel = nagios.StateUNKNOWNLabel
		stateExitCode = nagios.StateUNKNOWNExitCode
	}

	return nagios.ServiceState{
		Label:    stateLabel,
		ExitCode: stateExitCode,
	}

}

// UntilExpiration evaluates the given domain metadata and returns the number
// of days until the domain expires. If already expired, a negative number is
// returned indicating how many days the domain is past expiration.
//
// An error is returned if the pointer to the given domain metadata is nil.
func UntilExpiration(d *Domain) (int, error) {
	if d == nil {
		return 0, fmt.Errorf(
			"func UntilExpiration: unable to determine days until expiration: %w",
			ErrMissingValue,
		)
	}

	timeRemaining := time.Until(d.ExpirationDate).Hours()

	// Toss remainder so that we only get the whole number of days
	daysRemaining := int(math.Trunc(timeRemaining / 24))

	return daysRemaining, nil
}

// SinceUpdate evaluates the given domain metadata and returns the number of
// days since the domain metadata was last updated.
//
// An error is returned if the pointer to the given domain metadata is nil.
func SinceUpdate(d *Domain) (int, error) {
	if d == nil {
		return 0, fmt.Errorf(
			"func SinceUpdate: unable to determine days since last update: %w",
			ErrMissingValue,
		)
	}

	timeElapsed := time.Since(d.UpdatedDate).Hours()

	// Toss remainder so that we only get the whole number of days
	daysSince := int(math.Trunc(timeElapsed / 24))

	return daysSince, nil
}

// SinceCreation evaluates the given domain metadata and returns the number of
// days since the domain metadata was first created.
//
// An error is returned if the pointer to the given domain metadata is nil.
func SinceCreation(d *Domain) (int, error) {
	if d == nil {
		return 0, fmt.Errorf(
			"func SinceCreation: unable to determine days since creation: %w",
			ErrMissingValue,
		)
	}

	timeElapsed := time.Since(d.CreatedDate).Hours()

	// Toss remainder so that we only get the whole number of days
	daysSince := int(math.Trunc(timeElapsed / 24))

	return daysSince, nil
}

// FormattedExpiration receives a Time value and converts it to a string
// representing the largest useful whole units of time in days and hours. For
// example, if a domain has 1 year, 2 days and 3 hours remaining until
// expiration, this function will return the string '367d 3h remaining', but
// if only 3 hours remain then '3h remaining' will be returned. If a domain
// registration has expired, the 'ago' suffix will be used instead. For
// example, if a domain has expired 3 hours ago, '3h ago' will be returned.
func FormattedExpiration(expireTime time.Time) string {

	timeRemaining := time.Until(expireTime).Hours()

	var timeExpired bool
	var formattedTimeRemainingStr string
	var daysRemainingStr string
	var hoursRemainingStr string

	// Flip sign back to positive, note that expiraton has been reached for
	// later use.
	if timeRemaining < 0 {
		timeExpired = true
		timeRemaining *= -1
	}

	// Toss remainder so that we only get the whole number of days
	daysRemaining := math.Trunc(timeRemaining / 24)

	if daysRemaining > 0 {
		daysRemainingStr = fmt.Sprintf("%dd", int64(daysRemaining))
	}

	// Multiply the whole number of days by 24 to get the hours value, then
	// subtract from the original number of hours until expiration to get the
	// number of hours leftover from the days calculation.
	hoursRemaining := math.Trunc(timeRemaining - (daysRemaining * 24))

	hoursRemainingStr = fmt.Sprintf("%dh", int64(hoursRemaining))

	formattedTimeRemainingStr = strings.Join([]string{
		daysRemainingStr, hoursRemainingStr}, " ")

	switch {
	case !timeExpired:
		formattedTimeRemainingStr = strings.Join([]string{formattedTimeRemainingStr, "remaining"}, " ")
	case timeExpired:
		formattedTimeRemainingStr = strings.Join([]string{formattedTimeRemainingStr, "ago"}, " ")
	}

	return formattedTimeRemainingStr

}

// humanizeDays takes a target time and relative to a "base" time (usually
// time.Now()) returns a human readable phrase describing the difference in
// days.
func humanizeDays(target, base time.Time) string {
	// Normalize both times to midnight in their local timezones to get exact
	// calendar days.
	targetMidnight := time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location())
	baseMidnight := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, base.Location())

	days := int(math.Round(targetMidnight.Sub(baseMidnight).Hours() / 24))

	switch {
	case days == 0:
		return "today"

	case days == 1:
		return "tomorrow"

	case days == -1:
		return "yesterday"

	case days > 1:
		return fmt.Sprintf("in %d days", days)

	default:
		return fmt.Sprintf("%d days ago", int(math.Abs(float64(days))))
	}
}
