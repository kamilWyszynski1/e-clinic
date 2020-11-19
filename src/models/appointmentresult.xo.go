// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// AppointmentResult represents a row from 'public.appointment_result'.
var (
	AppointmentResultFields  = ` id, appointment, comment, prescription `
	AppointmentResultColumns = []string{"id", "appointment", "comment", "prescription"}
)

type AppointmentResult struct {
	ID           uuid.UUID `json:"id,omitempty"`           // id
	Appointment  uuid.UUID `json:"appointment,omitempty"`  // appointment
	Comment      string    `json:"comment,omitempty"`      // comment
	Prescription string    `json:"prescription,omitempty"` // prescription

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the AppointmentResult exists in the database.
func (ar *AppointmentResult) Exists() bool {
	return ar._exists
}

// Only for tests usage!
func (ar *AppointmentResult) TestOnly_SetExists() {
	ar._exists = true
}

// Only for tests usage!
func (ar *AppointmentResult) TestOnly_SetDeleted() {
	ar._deleted = true
}

// Deleted provides information if the AppointmentResult has been deleted from the database.
func (ar *AppointmentResult) Deleted() bool {
	return ar._deleted
}

// Insert inserts the AppointmentResult to the database.
func (ar *AppointmentResult) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if ar._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.appointment_result (` +
		`id, appointment, comment, prescription` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, ar.ID, ar.Appointment, ar.Comment, ar.Prescription)
	err = db.QueryRow(sqlstr, ar.ID, ar.Appointment, ar.Comment, ar.Prescription).Scan(&ar.ID)
	if err != nil {
		return err
	}

	// set existence
	ar._exists = true

	return nil
}

// Update updates the AppointmentResult in the database.
func (ar *AppointmentResult) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ar._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if ar._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.appointment_result SET (` +
		`appointment, comment, prescription` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	// run query
	XOLog(sqlstr, ar.Appointment, ar.Comment, ar.Prescription, ar.ID)
	_, err = db.Exec(sqlstr, ar.Appointment, ar.Comment, ar.Prescription, ar.ID)
	return err
}

// Save saves the AppointmentResult to the database.
func (ar *AppointmentResult) Save(db XODB) error {
	if ar.Exists() {
		return ar.Update(db)
	}

	return ar.Insert(db)
}

// Upsert performs an upsert for AppointmentResult.
//
// NOTE: PostgreSQL 9.5+ only
func (ar *AppointmentResult) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if ar._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.appointment_result (` +
		`id, appointment, comment, prescription` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, appointment, comment, prescription` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.appointment, EXCLUDED.comment, EXCLUDED.prescription` +
		`)`

	// run query
	XOLog(sqlstr, ar.ID, ar.Appointment, ar.Comment, ar.Prescription)
	_, err = db.Exec(sqlstr, ar.ID, ar.Appointment, ar.Comment, ar.Prescription)
	if err != nil {
		return err
	}

	// set existence
	ar._exists = true

	return nil
}

// Delete deletes the AppointmentResult from the database.
func (ar *AppointmentResult) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ar._exists {
		return nil
	}

	// if deleted, bail
	if ar._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.appointment_result WHERE id = $1`

	// run query
	XOLog(sqlstr, ar.ID)
	_, err = db.Exec(sqlstr, ar.ID)
	if err != nil {
		return err
	}

	// set deleted
	ar._deleted = true

	return nil
}

// AppointmentByAppointmentResultAppointmentFkey returns the Appointment associated with the AppointmentResult's Appointment (appointment).
//
// Generated from foreign key 'appointment_result_appointment_fkey'.
func (ar *AppointmentResult) AppointmentByAppointmentResultAppointmentFkey(db XODB) (*Appointment, error) {

	return AppointmentByID(db, ar.Appointment)

}

// AppointmentResultByID retrieves a row from 'public.appointment_result' as a AppointmentResult.
//
// Generated from index 'appointment_result_pkey'.
func AppointmentResultByID(db XODB, id uuid.UUID) (*AppointmentResult, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, appointment, comment, prescription ` +
		`FROM public.appointment_result ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	ar := AppointmentResult{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&ar.ID, &ar.Appointment, &ar.Comment, &ar.Prescription)
	if err != nil {
		return nil, err
	}

	return &ar, nil
}
