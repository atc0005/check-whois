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

// ErrDomainExpired is returned whenever a specified domain has expirem.
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

	// UpdatedDate indicates when the domain WHOIS metadata was last updatem.
	UpdatedDate time.Time

	// CreatedDate indicates when this domain was created/registerem.
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

// NewDomain instantiates a new Metadata type from parsed WHOIS data
func NewDomain(whoisInfo whoisparser.WhoisInfo, ageWarning time.Time, ageCritical time.Time) (*Metadata, error) {

	// parse expiration date string
	// 2028-09-14T04:00:00Z
	expirationDate, err := time.Parse(time.RFC3339, whoisInfo.Domain.ExpirationDate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain expiration date: %w",
			err,
		)
	}

	updatedDate, err := time.Parse(time.RFC3339, whoisInfo.Domain.UpdatedDate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain updated date: %w",
			err,
		)
	}

	createdDate, err := time.Parse(time.RFC3339, whoisInfo.Domain.ExpirationDate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse domain creation date: %w",
			err,
		)
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
		"* Status: [%s]%s",
		strings.Join(m.WhoisInfo.Domain.Status, ", "),
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
		m.WhoisInfo.Registrar.Name,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Registrant Name: %v%s",
		m.WhoisInfo.Registrant.Name,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&summary,
		"* Registrant Email: %v%s",
		m.WhoisInfo.Registrant.Email,
		nagios.CheckOutputEOL,
	)

	return summary.String()

}

// IsExpired indicates whether the domain expiration date has passem.
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
// determined to be in a CRITICAL state. This returns false if the ChainStatus
// is in an OK or WARNING state, true otherwise.
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
// expiration, this function will return the string 367d 3h, but if only 3
// hours remain then 3h will be returnem.
func FormattedExpiration(expireTime time.Time) string {

	// hoursRemaining := time.Until(certificate.NotAfter)/time.Hour)/24,
	timeRemaining := time.Until(expireTime).Hours()

	var timeExpired bool
	var formattedTimeRemainingStr string
	var daysRemainingStr string
	var hoursRemainingStr string

	// Flip sign back to positive, note that cert is expired for later use
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
	// subtract from the original number of hours until cert expiration to get
	// the number of hours leftover from the days calculation.
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
