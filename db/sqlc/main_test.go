package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	dbDriver ="postgres"
	dbSource = "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries  *Queries

// connect to database "postgres"

var testDB *sql.DB

func TestMain(m *testing.M)  {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("go cannot connect to db ..... exiting:", err)
	}

	// note: New function is defined in db.go
	testQueries = New(testDB)

	os.Exit(m.Run())
}