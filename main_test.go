package main

//https://stackoverflow.com/search?q=manifoldco%2Fpromptui+testing

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
	"work/models"

	"github.com/manifoldco/promptui"
)

func pad(siz int, buf *bytes.Buffer) {
	pu := make([]byte, 4096-siz)
	for i := 0; i < 4096-siz; i++ {
		pu[i] = 97
	}
	buf.Write(pu)
}

func TestMock(t *testing.T) {

	i1 := "N\n"
	i2 := "Y\n"

	b := bytes.NewBuffer([]byte(i1))
	pad(len(i1), b)
	reader := ioutil.NopCloser(
		b,
	)
	b.WriteString(i2)
	pad(len(i2), b)

	p := promptui.Prompt{
		Stdin: reader,
	}

	response := mock(p)

	if !strings.EqualFold(response, "NY") {
		t.Errorf("nope!")
		t.Errorf(response)
	}
}

func TestAddToDos(t *testing.T) {
	todos := []models.ToDo{}

	i1 := "to do 1\n"
	i2 := "to do 2\n"
	i3 := "\n"

	b := bytes.NewBuffer([]byte(i1))
	b.WriteString(i2)
	b.WriteString(i3)

	app1 := &app{
		todos: todos,
		r:     bufio.NewReader(b),
	}

	if len(app1.todos) != 0 {
		t.Errorf("to dos not set up correctly")
	}

	addToDos(app1)

	if len(app1.todos) != 2 {
		t.Errorf("to dos not added")
	}

	if "to do 1" != app1.todos[0].Title {
		t.Errorf("title not parsed correctly: %#v", app1.todos[0].Title)
	}

	//TEST WINDOWS
	i1 = "to do 3\r\n"
	i3 = "\r\n"

	b = bytes.NewBuffer([]byte(i1))
	b.WriteString(i3)

	app2 := &app{
		todos: todos,
		r:     bufio.NewReader(b),
	}

	addToDos(app2)

	if "to do 3" != app2.todos[0].Title {
		t.Errorf("title not parsed correctly: %#v", app2.todos[0].Title)
	}

}

func TestDistracted(t *testing.T) {

	// db_dir, err := os.MkdirTemp("testdata", "*")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	todos := []models.ToDo{}
	todo := models.ToDo{
		Title: "something",
	}
	todos = append(todos, todo)

	i1 := "a small thing\n"
	buf := bytes.NewBuffer([]byte(i1))

	app := &app{
		todos: todos,
		// db_dir: db_dir,
		chosen: todo,
		r:      bufio.NewReader(buf),
		sleep:  func(time.Duration) {},
	}

	if app.chosen.Title != "something" {
		t.Errorf("not set up correctly")
	}

	distracted(app)

	if app.chosen.Smalltask != strings.Trim(i1, "\n") {
		t.Errorf("failed to assign small task")
	}

	// os.RemoveAll(db_dir)
}

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
		// pprompt: promptui.Prompt{
		// 	Stdin: TestStdout,
		// },
	}

	app.db = models.ConnectDB(app.db_dir)
	if err := app.db.Ping(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println("This will now rm -fr " + db_dir)
	// _, _ = app.r.ReadString('\n')
	run(app)

	// bufio.NewWriter(os.Stdin).WriteString("Y\n")

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
