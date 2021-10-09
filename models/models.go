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
	Id    int
	Title string
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
	statement, err := c.DB.Prepare("INSERT INTO feelings (Title) VALUES (?)")
	if err != nil {
		println(fmt.Errorf(err.Error()))
	}
	statement.Exec(title)
}
