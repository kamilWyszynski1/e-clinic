// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// Substance represents a row from 'public.substance'.
var (
	SubstanceFields  = ` id, name `
	SubstanceColumns = []string{"id", "name"}
)

type Substance struct {
	ID   uuid.UUID `json:"id,omitempty"`   // id
	Name string    `json:"name,omitempty"` // name

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Substance exists in the database.
func (s *Substance) Exists() bool {
	return s._exists
}

// Only for tests usage!
func (s *Substance) TestOnly_SetExists() {
	s._exists = true
}

// Only for tests usage!
func (s *Substance) TestOnly_SetDeleted() {
	s._deleted = true
}

// Deleted provides information if the Substance has been deleted from the database.
func (s *Substance) Deleted() bool {
	return s._deleted
}

// Insert inserts the Substance to the database.
func (s *Substance) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if s._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.substance (` +
		`id, name` +
		`) VALUES (` +
		`$1, $2` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, s.ID, s.Name)
	err = db.QueryRow(sqlstr, s.ID, s.Name).Scan(&s.ID)
	if err != nil {
		return err
	}

	// set existence
	s._exists = true

	return nil
}

// Update updates the Substance in the database.
func (s *Substance) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if s._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.substance SET (` +
		`name` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	// run query
	XOLog(sqlstr, s.Name, s.ID)
	_, err = db.Exec(sqlstr, s.Name, s.ID)
	return err
}

// Save saves the Substance to the database.
func (s *Substance) Save(db XODB) error {
	if s.Exists() {
		return s.Update(db)
	}

	return s.Insert(db)
}

// Upsert performs an upsert for Substance.
//
// NOTE: PostgreSQL 9.5+ only
func (s *Substance) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if s._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.substance (` +
		`id, name` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name` +
		`)`

	// run query
	XOLog(sqlstr, s.ID, s.Name)
	_, err = db.Exec(sqlstr, s.ID, s.Name)
	if err != nil {
		return err
	}

	// set existence
	s._exists = true

	return nil
}

// Delete deletes the Substance from the database.
func (s *Substance) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !s._exists {
		return nil
	}

	// if deleted, bail
	if s._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.substance WHERE id = $1`

	// run query
	XOLog(sqlstr, s.ID)
	_, err = db.Exec(sqlstr, s.ID)
	if err != nil {
		return err
	}

	// set deleted
	s._deleted = true

	return nil
}

// SubstanceByID retrieves a row from 'public.substance' as a Substance.
//
// Generated from index 'substance_pkey'.
func SubstanceByID(db XODB, id uuid.UUID) (*Substance, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name ` +
		`FROM public.substance ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	s := Substance{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&s.ID, &s.Name)
	if err != nil {
		return nil, err
	}

	return &s, nil
}