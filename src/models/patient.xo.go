// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// Patient represents a row from 'public.patient'.
var (
	PatientFields  = ` id, name, surname, email, phone_number, age, gender `
	PatientColumns = []string{"id", "name", "surname", "email", "phone_number", "age", "gender"}
)

type Patient struct {
	ID          uuid.UUID  `json:"id,omitempty"`           // id
	Name        string     `json:"name,omitempty"`         // name
	Surname     string     `json:"surname,omitempty"`      // surname
	Email       string     `json:"email,omitempty"`        // email
	PhoneNumber string     `json:"phone_number,omitempty"` // phone_number
	Age         int        `json:"age,omitempty"`          // age
	Gender      Genderenum `json:"gender,omitempty"`       // gender

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Patient exists in the database.
func (p *Patient) Exists() bool {
	return p._exists
}

// Only for tests usage!
func (p *Patient) TestOnly_SetExists() {
	p._exists = true
}

// Only for tests usage!
func (p *Patient) TestOnly_SetDeleted() {
	p._deleted = true
}

// Deleted provides information if the Patient has been deleted from the database.
func (p *Patient) Deleted() bool {
	return p._deleted
}

// Insert inserts the Patient to the database.
func (p *Patient) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if p._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.patient (` +
		`id, name, surname, email, phone_number, age, gender` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, p.ID, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender)
	err = db.QueryRow(sqlstr, p.ID, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender).Scan(&p.ID)
	if err != nil {
		return err
	}

	// set existence
	p._exists = true

	return nil
}

// Update updates the Patient in the database.
func (p *Patient) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !p._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if p._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.patient SET (` +
		`name, surname, email, phone_number, age, gender` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6` +
		`) WHERE id = $7`

	// run query
	XOLog(sqlstr, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender, p.ID)
	_, err = db.Exec(sqlstr, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender, p.ID)
	return err
}

// Save saves the Patient to the database.
func (p *Patient) Save(db XODB) error {
	if p.Exists() {
		return p.Update(db)
	}

	return p.Insert(db)
}

// Upsert performs an upsert for Patient.
//
// NOTE: PostgreSQL 9.5+ only
func (p *Patient) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if p._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.patient (` +
		`id, name, surname, email, phone_number, age, gender` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, surname, email, phone_number, age, gender` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.surname, EXCLUDED.email, EXCLUDED.phone_number, EXCLUDED.age, EXCLUDED.gender` +
		`)`

	// run query
	XOLog(sqlstr, p.ID, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender)
	_, err = db.Exec(sqlstr, p.ID, p.Name, p.Surname, p.Email, p.PhoneNumber, p.Age, p.Gender)
	if err != nil {
		return err
	}

	// set existence
	p._exists = true

	return nil
}

// Delete deletes the Patient from the database.
func (p *Patient) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !p._exists {
		return nil
	}

	// if deleted, bail
	if p._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.patient WHERE id = $1`

	// run query
	XOLog(sqlstr, p.ID)
	_, err = db.Exec(sqlstr, p.ID)
	if err != nil {
		return err
	}

	// set deleted
	p._deleted = true

	return nil
}

// PatientByID retrieves a row from 'public.patient' as a Patient.
//
// Generated from index 'patient_pkey'.
func PatientByID(db XODB, id uuid.UUID) (*Patient, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, surname, email, phone_number, age, gender ` +
		`FROM public.patient ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	p := Patient{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&p.ID, &p.Name, &p.Surname, &p.Email, &p.PhoneNumber, &p.Age, &p.Gender)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
