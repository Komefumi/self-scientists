package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	_ "github.com/lib/pq"
)

var envValueRegexp *regexp.Regexp = regexp.MustCompile("(\\w+)=(.+)")
var DB *sql.DB

func init() {
	setEnvironmentVariables()
	setupDB()
}

func setEnvironmentVariables() {
	res, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal(err)
	}
	matches := envValueRegexp.FindAllStringSubmatch(string(res), -1)
	for _, v := range matches {
		os.Setenv(v[1], v[2])
	}
	// fmt.Println(os.Getenv("SERVER_PORT"))
	// fmt.Println(os.Getenv("DB_DSN"))
}

func setupDB() {
	// _, err := sql.Open("postgres", "dbname=self-scientists; sslmode=disable;")
	var err error
	DB, err = sql.Open(os.Getenv("DB_TYPE"), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection Established!")
	// fmt.Println(DB)
}

func main() {
	defer DB.Close()
	var col string
	row := DB.QueryRow("SELECT 10;")
	err := row.Scan(&col)
	if err != nil {
		log.Fatal(err)
	}
}
