/*
 * Copyright 2014-2021 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Go module for domain whois information parsing
 * https://www.likexian.com/
 */

package whoisparser

import (
	"errors"
	"strings"
)

var (
	// ErrDomainNotFound domain is not found
	ErrDomainNotFound = errors.New("whoisparser: domain is not found")
	// ErrReservedDomain domain is reserved
	ErrReservedDomain = errors.New("whoisparser: domain is reserved to register")
	// ErrPremiumDomain domain is available to register at premium price
	ErrPremiumDomain = errors.New("whoisparser: domain is available at premium price")
	// ErrBlockedDomain domain is blocked due to brand protection
	ErrBlockedDomain = errors.New("whoisparser: domain is blocked due to brand protection")
	// ErrDomainDataInvalid domain whois data is invalid
	ErrDomainDataInvalid = errors.New("whoisparser: domain whois data is invalid")
	// ErrDomainLimitExceed domain whois query is limited
	ErrDomainLimitExceed = errors.New("whoisparser: domain whois query limit exceeded")
)

// getDomainErrorType returns error type of domain data
func getDomainErrorType(data string) error {
	switch {
	case isDomainNotFound(data):
		return ErrDomainNotFound
	case isBlockedDomain(data):
		return ErrBlockedDomain
	case isPremiumDomain(data):
		return ErrPremiumDomain
	case isReservedDomain(data):
		return ErrReservedDomain
	case isLimitExceeded(data):
		return ErrDomainLimitExceed
	default:
		return ErrDomainDataInvalid
	}
}

// isDomainNotFound returns if domain is not found
func isDomainNotFound(data string) bool {
	notFoundKeys := []string{
		"is free",
		"no found",
		"no match",
		"not found",
		"not match",
		"no data found",
		"no entries found",
		"no matching record",
		"not registered",
		"not been registered",
		"object does not exist",
	}

	return containsIn(strings.ToLower(data), notFoundKeys)
}

// isReservedDomain returns if domain is reserved
func isReservedDomain(data string) bool {
	reservedKeys := []string{
		"reserved domain name",
		"reserved by the registry",
		"can not be registered online",
	}

	return containsIn(strings.ToLower(data), reservedKeys)
}

// isPremiumDomain returns if domain is available to register at premium price
func isPremiumDomain(data string) bool {
	premiumKeys := []string{
		"premium domain is available for purchase",
		"platinum domain is available for purchase",
	}

	return containsIn(strings.ToLower(data), premiumKeys)
}

// isBlockedDomain returns if domain is blocked due to brand protection
func isBlockedDomain(data string) bool {
	blockedKeys := []string{
		// Donuts DPML
		"dpml brand protection",
		// Uniregistry Uni EPS
		"subscribes to the uni eps",
		// Gandi AdultBlock
		"subscribes to the adultblock",
	}

	return containsIn(strings.ToLower(data), blockedKeys)
}

// isLimitExceeded returns if domain whois query is limited
func isLimitExceeded(data string) bool {
	limitExceedKeys := []string{
		"limit exceeded",
	}

	return containsIn(strings.ToLower(data), limitExceedKeys)
}
