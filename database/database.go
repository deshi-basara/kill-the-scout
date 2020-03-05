package database

import (
	"log"
	"time"

	"github.com/asdine/storm/v3"
)

// Expose represents one immobilienscout expose
type Expose struct {
	ID        string `storm:"unique"` // primary key
	URL       string
	Title     string
	Address   string
	Rent      string
	Size      string
	Rooms     string
	CreatedAt time.Time
}

// Setup local database
func Setup() *storm.DB {
	db, err := storm.Open("immo.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// SaveExpose inserts the handed Expose into embedded database
func SaveExpose(db *storm.DB, expose Expose) error {
	err := db.Save(&expose)

	return err
}

// GetExpose tries to find and return an Expose identified by exposeID
func GetExpose(db *storm.DB, exposeID string) (Expose, error) {
	var expose Expose
	err := db.One("ID", exposeID, &expose)

	return expose, err
}
