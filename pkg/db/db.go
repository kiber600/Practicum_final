package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(32) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCAHAR(7) NOT NULL DEFAULT ""
	);`

var db *sql.DB

func CheckDbFile() error {
	dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)
	var install bool
	if err != nil {
		install = true
	}
	if install == true {
		err = createDb()
		if err != nil {
			return err
		}
	}

	if install == false {
		err = openDb()
		if err != nil {
			return err
		}
	}
	return nil
}

func createDb() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return err
	}

	db, err = sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		return err
	}
	_, err = db.Exec(schema)
	if err != nil {
		log.Printf("%s\n", err.Error())
		return err
	}
	return nil
}

func openDb() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return err
	}

	db, err = sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		return err
	}

	return nil
}

func CloseDB() error {
	return db.Close()
}
