// Package models contains the types for schema 'public'.
package models

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	uuid "github.com/satori/go.uuid"
)

// Composition represents a row from 'public.composition'.
var (
	CompositionFields  = ` id, drug, substance `
	CompositionColumns = []string{"id", "drug", "substance"}
)

type Composition struct {
	ID        uuid.UUID `json:"id,omitempty"`        // id
	Drug      int       `json:"drug,omitempty"`      // drug
	Substance uuid.UUID `json:"substance,omitempty"` // substance

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Composition exists in the database.
func (c *Composition) Exists() bool {
	return c._exists
}

// Only for tests usage!
func (c *Composition) TestOnly_SetExists() {
	c._exists = true
}

// Only for tests usage!
func (c *Composition) TestOnly_SetDeleted() {
	c._deleted = true
}

// Deleted provides information if the Composition has been deleted from the database.
func (c *Composition) Deleted() bool {
	return c._deleted
}

// Insert inserts the Composition to the database.
func (c *Composition) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO public.composition (` +
		`id, drug, substance` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`

	// run query
	XOLog(sqlstr, c.ID, c.Drug, c.Substance)
	err = db.QueryRow(sqlstr, c.ID, c.Drug, c.Substance).Scan(&c.ID)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Update updates the Composition in the database.
func (c *Composition) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if c._deleted {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE public.composition SET (` +
		`drug, substance` +
		`) = ( ` +
		`$1, $2` +
		`) WHERE id = $3`

	// run query
	XOLog(sqlstr, c.Drug, c.Substance, c.ID)
	_, err = db.Exec(sqlstr, c.Drug, c.Substance, c.ID)
	return err
}

// Save saves the Composition to the database.
func (c *Composition) Save(db XODB) error {
	if c.Exists() {
		return c.Update(db)
	}

	return c.Insert(db)
}

// Upsert performs an upsert for Composition.
//
// NOTE: PostgreSQL 9.5+ only
func (c *Composition) Upsert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql query
	const sqlstr = `INSERT INTO public.composition (` +
		`id, drug, substance` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, drug, substance` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.drug, EXCLUDED.substance` +
		`)`

	// run query
	XOLog(sqlstr, c.ID, c.Drug, c.Substance)
	_, err = db.Exec(sqlstr, c.ID, c.Drug, c.Substance)
	if err != nil {
		return err
	}

	// set existence
	c._exists = true

	return nil
}

// Delete deletes the Composition from the database.
func (c *Composition) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return nil
	}

	// if deleted, bail
	if c._deleted {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM public.composition WHERE id = $1`

	// run query
	XOLog(sqlstr, c.ID)
	_, err = db.Exec(sqlstr, c.ID)
	if err != nil {
		return err
	}

	// set deleted
	c._deleted = true

	return nil
}

// DrugByCompositionDrugFkey returns the Drug associated with the Composition's Drug (drug).
//
// Generated from foreign key 'composition_drug_fkey'.
func (c *Composition) DrugByCompositionDrugFkey(db XODB) (*Drug, error) {

	return DrugByID(db, c.Drug)

}

// SubstanceByCompositionSubstanceFkey returns the Substance associated with the Composition's Substance (substance).
//
// Generated from foreign key 'composition_substance_fkey'.
func (c *Composition) SubstanceByCompositionSubstanceFkey(db XODB) (*Substance, error) {

	return SubstanceByID(db, c.Substance)

}

// CompositionByID retrieves a row from 'public.composition' as a Composition.
//
// Generated from index 'composition_pkey'.
func CompositionByID(db XODB, id uuid.UUID) (*Composition, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, drug, substance ` +
		`FROM public.composition ` +
		`WHERE id = $1`

	// run query
	XOLog(sqlstr, id)
	c := Composition{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&c.ID, &c.Drug, &c.Substance)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
