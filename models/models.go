package models

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func ConnectDB(db_dir string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath.Join(db_dir, "work_motivated.db"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type Feeling struct {
	id          int
	Type        string
	Description string
}

type FeelingModel struct {
	DB *sql.DB
}

type ToDo struct {
	Id        int
	Title     string
	Smalltask string
}

type ToDoModel struct {
	DB *sql.DB
}

func (c *FeelingModel) CreateTable() {
	statement, err := c.DB.Prepare("CREATE TABLE IF NOT EXISTS feelings (id INTEGER PRIMARY KEY, Type TEXT, Description TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statement.Exec()
	if err != nil {
		fmt.Printf("%#v", r)
		log.Fatal(err)
	}
}

func (c *FeelingModel) Insert(dtype string, description string) {
	if strings.EqualFold("\n", description) {
		return
	}
	c.CreateTable()
	statement, err := c.DB.Prepare("INSERT INTO feelings (Type, Description) VALUES (?, ?)")
	if err != nil {
		println(fmt.Errorf(err.Error()))
	}
	statement.Exec(dtype, description)
}

//BUG: no such table: feelings, no rows in result set ERRORS SHOWING WHEN TRYING TO INSERT INTO A NON-EXISTANT TABLE
func (c *FeelingModel) GetRandom(dtype string) Feeling {

	feeling := Feeling{}
	where := ""
	if dtype != "" {
		where = "WHERE Type = '" + dtype + "'"
	}
	query := "SELECT Type, Description FROM feelings " + where + " ORDER BY RANDOM() LIMIT 1"
	err := c.DB.QueryRow(query).Scan(&feeling.Type, &feeling.Description)
	if err != nil {
		//CHANGE THIS TO NOT SHOW ERROR sql: no rows in result set
		fmt.Println(fmt.Errorf(err.Error()))
	}
	return feeling
}

func (c *ToDoModel) CreateTable() {
	statement, err := c.DB.Prepare("CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, Title TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	r, err := statement.Exec()
	if err != nil {
		fmt.Printf("%#v", r)
		log.Fatal(err)
	}
}

func (c *ToDoModel) Insert(title string) {
	c.CreateTable()
	statement, err := c.DB.Prepare("INSERT INTO todos (Title) VALUES (?)")
	if err != nil {
		println(fmt.Errorf(err.Error()))
	}
	statement.Exec(title)
}
