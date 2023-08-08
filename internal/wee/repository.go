// Module repository provides functions to store wee Urls in a lookup table.
package wee

import (
	"fmt"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

// Database operations for Wee Records
// This implementation uses an SQLite DB
const (
	tableName = "weeRecords"

	// Using the sqlite column types (instead of the go ones) for compatibility
	tableDefn = `"Tag" TEXT UNIQUE ON CONFLICT FAIL, "Url" TEXT, "Token" TEXT`
)

// Record defines and works with the Wee record entries.
// Storage of the records in a table is abstracted to the db module.

type Record struct {
	Tag		string
	Url		string
	Token	string
}

// Repository bundles the config of the table and a handle to it
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

// Connect opens the database and creates the Url table if it doesn't yet exist.
func (r *Repository) Connect() error {

	var err error
	var db *sql.DB
	
	// Create the directory for the repository, if necessary -- say, if one a volume
	dir := filepath.Dir(r.source)
	if dir != "." {
		err = os.MkdirAll(dir, 0750)
		if err != nil {
			r.logger.Printf("Error creating dir for %s: %v\n", r.source, err)
			return err
		}
	}
	
	// Connect to specified service
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

// Disconnect closes the connection.
func (r *Repository) Disconnect() error {
	// TBD anything else to do?
	return r.db.Close()
}

// create the table used to persist Urls.
// TBD locking should be added to avoid concurrency problems
func (r *Repository) create() error {

	cmd := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s ( %s );`, tableName, tableDefn)
	r.logger.Printf("Exec SQL: %s\n", cmd)

	_, err := r.db.Exec(cmd)
	if err != nil {
		r.logger.Printf("ERROR on SQL table creation %s, %v\n", tableName, err)
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
// If error is nil then there IS a match,
// ErrNoRows is returned on the no-match condition, and other errors are possible, too;
// Record will be nil on any error.
func (r *Repository) find(weeUrl string) (*Record, error) {

	var err error
	var rec Record

	// it's a unique tag so there should be no more than a single row 
	err = r.db.QueryRow("SELECT * FROM " + tableName + " WHERE Tag = ?", weeUrl).Scan(&rec.Tag, &rec.Url, &rec.Token)

	if err != nil {
		return nil, err
	}

	return &rec, err
}

// remove deletes the Record having the specified token.
// TBD TestRemove has found that an error is NOT reported when removing a record
//     with a non-existent token.
func (r *Repository) remove(token string) error {

	stmt, err := r.db.Prepare("DELETE FROM " + tableName + " WHERE Token = ?")
	_, err = stmt.Exec(token)
	if err != nil {
		return err
	}

	return nil
}

