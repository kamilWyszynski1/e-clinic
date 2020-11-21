// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql/driver"
	"errors"
)

// Apoitntmentstateenum is the 'apoitntmentstateenum' enum type from schema 'public'.
type Apoitntmentstateenum string

const (
	// ApoitntmentstateenumCreated is the 'CREATED' Apoitntmentstateenum.
	ApoitntmentstateenumCreated = Apoitntmentstateenum("CREATED")

	// ApoitntmentstateenumAccepted is the 'ACCEPTED' Apoitntmentstateenum.
	ApoitntmentstateenumAccepted = Apoitntmentstateenum("ACCEPTED")

	// ApoitntmentstateenumToBePaid is the 'TO_BE_PAID' Apoitntmentstateenum.
	ApoitntmentstateenumToBePaid = Apoitntmentstateenum("TO_BE_PAID")

	// ApoitntmentstateenumOk is the 'OK' Apoitntmentstateenum.
	ApoitntmentstateenumOk = Apoitntmentstateenum("OK")

	// ApoitntmentstateenumFinished is the 'FINISHED' Apoitntmentstateenum.
	ApoitntmentstateenumFinished = Apoitntmentstateenum("FINISHED")

	// ApoitntmentstateenumRejected is the 'REJECTED' Apoitntmentstateenum.
	ApoitntmentstateenumRejected = Apoitntmentstateenum("REJECTED")

	// ApoitntmentstateenumNotPaid is the 'NOT_PAID' Apoitntmentstateenum.
	ApoitntmentstateenumNotPaid = Apoitntmentstateenum("NOT_PAID")
)

var ApoitntmentstateenumValues = []Apoitntmentstateenum{ApoitntmentstateenumCreated, ApoitntmentstateenumAccepted, ApoitntmentstateenumToBePaid, ApoitntmentstateenumOk, ApoitntmentstateenumFinished, ApoitntmentstateenumRejected, ApoitntmentstateenumNotPaid}

// String returns the string value of the Apoitntmentstateenum.
func (a Apoitntmentstateenum) String() string {
	return string(a)
}

// MarshalText marshals Apoitntmentstateenum into text.
func (a Apoitntmentstateenum) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// UnmarshalText unmarshals Apoitntmentstateenum from text.
func (a *Apoitntmentstateenum) UnmarshalText(text []byte) error {
	switch string(text) {
	case "CREATED":
		*a = ApoitntmentstateenumCreated

	case "ACCEPTED":
		*a = ApoitntmentstateenumAccepted

	case "TO_BE_PAID":
		*a = ApoitntmentstateenumToBePaid

	case "OK":
		*a = ApoitntmentstateenumOk

	case "FINISHED":
		*a = ApoitntmentstateenumFinished

	case "REJECTED":
		*a = ApoitntmentstateenumRejected

	case "NOT_PAID":
		*a = ApoitntmentstateenumNotPaid

	default:
		return errors.New("invalid Apoitntmentstateenum")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Apoitntmentstateenum.
func (a Apoitntmentstateenum) Value() (driver.Value, error) {
	return a.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Apoitntmentstateenum.
func (a *Apoitntmentstateenum) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Apoitntmentstateenum")
	}

	return a.UnmarshalText(buf)
}
