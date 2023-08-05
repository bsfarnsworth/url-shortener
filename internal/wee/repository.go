package wee

import (
	"fmt"
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

// Database operations for Wee Records
// This implementation uses an SQLite DB
const (
	tableName = "weeRecords"

	// TBD combine w/ Record struc
	tableDefn = `"Tag" TEXT UNIQUE ON CONFLICT FAIL, "Url" TEXT, "Token" TEXT`
)

// Record defines and works with the Wee record entries.
// Storage of the records in a table is abstracted to the db module.

type Record struct {
	Tag		string
	Url		string
	Token	string
}

type Repository struct {
	driver string					// expecting "sqlite3" or compatible
	source string					// per deployment, e.g. "./foo.db"
	logger *log.Logger
	db *sql.DB
}

func NewRepository(drv string, src string, log *log.Logger) *Repository {
	return &Repository{
		driver: drv,
		source: src,
		logger: log,
		db: nil}
}

func (r *Repository) Connect() error {

	// Connect to specified service
	var err error
	var db *sql.DB
	
	db, err = sql.Open(r.driver, r.source)
	r.db = db
	if err != nil {
		r.logger.Printf("Error while opening %s: %v\n", r.source, err)
		return err
	}
	
	// This will create the table only if it doesn't yet exist
	err = r.create()
	if err != nil {
		r.logger.Printf("Error while creating table %s: %v\n", r.source, err)
		return err
	}

	return nil
}

func (r *Repository) Disconnect() error {
	// TBD anything else to do?
	return r.db.Close()
}

func (r *Repository) create() error {

	cmd := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( %s );`, tableName, tableDefn)
	r.logger.Printf("Exec SQL: %s\n", cmd)

	_, err := r.db.Exec(cmd)
	if err != nil {
		r.logger.Printf("ERROR on SQL table creatiion %s, %v\n", tableName, err)
	}
	return err
}


// add inserts a Record into the table
func (r *Repository) add(rec *Record) error {
	stmt, _ := r.db.Prepare("INSERT INTO " + tableName + " VALUES (?, ?, ?)")
	_, err := stmt.Exec(rec.Tag, rec.Url, rec.Token)
	defer stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

// find locates an existing Record for the given weeUrl
func (r *Repository) find(weeUrl string) (*Record, error) {
	var rec Record
	stmt, _ := r.db.Prepare("SELECT * FROM " + tableName + " WHERE Tag = ?")
	_, err := stmt.Exec(weeUrl)
	//	defer row.Close()
	if err != nil {
		return nil, err
	}
	/*
	err = row.Err()
	if err != nil {
		return nil, err
	}
	var rec Record
	err = row.Scan(&rec.Tag, &rec.Url, &rec.Token)
	if err != nil {
	}
	*/
	return &rec, nil
}

func (r *Repository) remove(rec *Record) error {
	stmt, err := r.db.Prepare("DELETE FROM records WHERE tag = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(rec.Tag)
	return nil
}

