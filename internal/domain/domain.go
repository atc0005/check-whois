// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-whois
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	whoisparser "github.com/likexian/whois-parser"
)

// DomainDateLayout is the chosen date layout for displaying
// domain created/expiration date/time values across our application.
const DomainDateLayout string = "2006-01-02 15:04:05 -0700 MST"

// defaultWhoISPlaceholderValue is used as a fallback value for any not found
// in the WhoIS record data (e.g., registrant name or email).
const defaultWhoISPlaceholderValue string = "unspecified"

// ErrDomainExpired is returned whenever a specified domain has expired.
var ErrDomainExpired = errors.New("domain has expired")

// ErrDomainExpiring is returned whenever a specified domain is expiring.
var ErrDomainExpiring = errors.New("domain is expiring")

// Metadata represents the details for a specified domain, including the
// parsed WHOIS info, expiration (age) thresholds and parsed date values.
type Metadata struct {

	// Results from parsing the WHOIS info. Most fields are plain text values.
	WhoisInfo whoisparser.WhoisInfo

	// Name is the plaintext label for this domain.
	Name string

	// ExpirationDate indicates when this domain expires.
	ExpirationDate time.Time

	// UpdatedDate indicates when the domain WHOIS metadata was last updated.
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

// parseDateString attempts to parse a given date string using detailed
// formats first, then falls back to a basic format before finally giving up
// and returning an error if no formats match.
func parseDateString(dateString string) (time.Time, error) {
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

// NewDomain instantiates a new Metadata type from parsed WHOIS data.
func NewDomain(whoisInfo whoisparser.WhoisInfo, ageWarning time.Time, ageCritical time.Time) (*Metadata, error) {

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
		expirationDate, err = parseDateString(whoisInfo.Domain.ExpirationDate)
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
		updatedDate, err = parseDateString(whoisInfo.Domain.UpdatedDate)
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
		createdDate, err = parseDateString(whoisInfo.Domain.CreatedDate)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse domain creation date: %w",
				err,
			)
		}
	}

	d := Metadata{
		AgeWarningThreshold:  ageWarning,
		AgeCriticalThreshold: ageCritical,
		WhoisInfo:            whoisInfo,
		Name:                 whoisInfo.Domain.Domain,
		ExpirationDate:       expirationDate,
		UpdatedDate:          updatedDate,
		CreatedDate:          createdDate,
	}

	return &d, nil

}

// OneLineCheckSummary generates a one-line summary of the domain WHOIS check
// results for display and notification purposes.
func (m Metadata) OneLineCheckSummary() string {

	var summary string

	switch {
	case m.IsExpired():
		summary = fmt.Sprintf(
			"%s: %q domain registration EXPIRED %s%s",
			m.ServiceState().Label,
			m.Name,
			FormattedExpiration(m.ExpirationDate),
			nagios.CheckOutputEOL,
		)

	default:

		summary = fmt.Sprintf(
			"%s: %q domain registration has %s%s",
			m.ServiceState().Label,
			m.Name,
			FormattedExpiration(m.ExpirationDate),
			nagios.CheckOutputEOL,
		)

	}

	return summary

}

// Report provides an overview of domain details appropriate for display as
// the LongServiceOutput provided via the web UI or as email or Teams
// notifications.
func (m Metadata) Report() string {

	var summary strings.Builder

	fmt.Fprintf(
		&summary,
		"WHOIS metadata for %q domain:%s%s",
		m.Name,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Status: %s%s",
		domainStatus(m),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Creation Date: %v%s",
		m.CreatedDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Updated Date: %v%s",
		m.UpdatedDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Expiration Date: %v%s",
		m.ExpirationDate.Format(DomainDateLayout),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Registrar Name: %v%s",
		registrarName(m),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Registrant Name: %v%s",
		registrantName(m),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Registrant Email: %v%s",
		registrantEmail(m),
		nagios.CheckOutputEOL,
	)

	return summary.String()

}

// IsExpired indicates whether the domain expiration date has passed.
func (m Metadata) IsExpired() bool {
	return m.ExpirationDate.Before(time.Now())
}

// IsExpiring compares the domain's current expiration date against the
// provided CRITICAL and WARNING thresholds to determine if the domain is
// about to expire.
func (m Metadata) IsExpiring() bool {

	switch {
	case !m.IsExpired() && m.ExpirationDate.Before(m.AgeCriticalThreshold):
		return true
	case !m.IsExpired() && m.ExpirationDate.Before(m.AgeWarningThreshold):
		return true
	}

	return false

}

// IsWarningState indicates whether a domain's expiration date has been
// determined to be in a WARNING state. This returns false if the expiration
// date is in an OK or CRITICAL state, true otherwise.
func (m Metadata) IsWarningState() bool {
	if !m.IsExpired() &&
		m.ExpirationDate.Before(m.AgeWarningThreshold) &&
		!m.ExpirationDate.Before(m.AgeCriticalThreshold) {
		return true
	}

	return false

}

// IsCriticalState indicates whether a domain's expiration date has been
// determined to be in a CRITICAL state. This returns false if the expiration
// date is in an OK or WARNING state, true otherwise.
func (m Metadata) IsCriticalState() bool {
	if m.IsExpired() || m.ExpirationDate.Before(m.AgeCriticalThreshold) {
		return true
	}

	return false

}

// IsOKState indicates whether a domain's expiration date has been determined
// to be in an OK state, without expired or expiring domain registration.
func (m Metadata) IsOKState() bool {
	return !m.IsWarningState() && !m.IsCriticalState()
}

// ServiceState returns the appropriate Service Check Status label and exit
// code for the evaluated domain expiration metadata.
func (m Metadata) ServiceState() nagios.ServiceState {

	var stateLabel string
	var stateExitCode int

	switch {
	case m.IsCriticalState():
		stateLabel = nagios.StateCRITICALLabel
		stateExitCode = nagios.StateCRITICALExitCode
	case m.IsWarningState():
		stateLabel = nagios.StateWARNINGLabel
		stateExitCode = nagios.StateWARNINGExitCode
	case m.IsOKState():
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

// domainStatus provides the domain status value from the WhoIS record or the
// fallback/placeholder value for the field.
func domainStatus(m Metadata) string {
	if m.WhoisInfo.Domain != nil && len(m.WhoisInfo.Domain.Status) != 0 {
		return strings.Join(m.WhoisInfo.Domain.Status, ", ")
	}

	return defaultWhoISPlaceholderValue
}

// registrarName provides the registrar name value from the WhoIS record or
// the fallback/placeholder value for the field.
func registrarName(m Metadata) string {
	if m.WhoisInfo.Registrar != nil && m.WhoisInfo.Registrar.Name != "" {
		return m.WhoisInfo.Registrar.Name
	}

	return defaultWhoISPlaceholderValue
}

// registrantName provides the registrant name value from the WhoIS record or
// the fallback/placeholder value for the field.
func registrantName(m Metadata) string {
	if m.WhoisInfo.Registrant != nil && m.WhoisInfo.Registrant.Name != "" {
		return m.WhoisInfo.Registrant.Name
	}

	return defaultWhoISPlaceholderValue
}

// registrantEmail provides the registrant email value from the WhoIS record
// or the fallback/placeholder value for the field.
func registrantEmail(m Metadata) string {
	if m.WhoisInfo.Registrant != nil && m.WhoisInfo.Registrant.Email != "" {
		return m.WhoisInfo.Registrant.Email
	}

	return defaultWhoISPlaceholderValue
}
