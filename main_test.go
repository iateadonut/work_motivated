package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	"work/models"
)

func TestRun(t *testing.T) {

	db_dir, err := os.MkdirTemp("testdata", "*")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db_dir)

	app := &app{
		r:     bufio.NewReader(os.Stdin),
		todos: []models.ToDo{},
		//chosen
		//db:     *sql.DB,
		sleep:   func(d time.Duration) {},
		c_timer: false,
		db_dir:  db_dir,
	}

	app.db = models.ConnectDB(app.db_dir)
	if err := app.db.Ping(); err != nil {
		log.Fatal(err)
	}

	os.RemoveAll(db_dir)

}

func TestTimerFuncOnlyCreatesOneTimer(t *testing.T) {

	app := &app{
		//sleep: func(int) {},
		sleep: time.Sleep,
		//sleep_timer: sleep_timer,
		c_timer: false,
	}

	d := time.Duration(time.Millisecond * 50)

	if app.c_timer != false {
		t.Errorf("timer is not false")
	}

	go timer(app, d)

	time.Sleep(time.Millisecond * 10)
	if app.c_timer != true {
		t.Errorf("timer does not run")
	}

	time.Sleep(time.Millisecond * 50)

	if app.c_timer != false {
		t.Errorf("timer is not destroyed")
	}

	//make sure we're set up to run the next test
	if app.c_timer != false {
		t.Errorf("timer is not false")
	}
	//this timer should make app.c_timer true for 400 ms
	go timer(app, time.Millisecond*400)
	//this timer would make app.c_timer false after 50ms if the first one did not already exist, so this should not have any effect on app.c_timer
	time.Sleep(time.Millisecond * 10)
	if app.c_timer != true {
		t.Errorf("timer does not run")
	}
	go timer(app, d)

	time.Sleep(time.Millisecond * 60)

	if app.c_timer == false {
		t.Errorf("timer overwrites old timer")
	}

}
