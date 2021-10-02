package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"work/models"

	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
)

func choose_one(reader bufio.Reader, todos []models.ToDo) {

	//if input was "ok", meaning you got distracted
	//are you working on another, more important task now?
	//are you bored?  - what's the smallest possible thing you can do to get started?
	//let's try to see the big picture
	//
	//skill acquisition? //so good they can't ignore you
	//flow?
	//produce results?

	// fmt.Println("Choose one thing that you should get started on:")
	// to_do, err := reader.ReadString('\n')
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("%v", todos)

	items := make([]string, 0, len(todos))
	for _, t := range todos {
		items = append(items, strings.Trim(t.Title, "\n"))
	}

	prompt := promptui.Select{
		Label: "Choose one thing that you should get started on",
		Items: items,
	}

	_, to_do, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Println()

	fmt.Println("Repeat to yourself:")
	fmt.Println("I will work with intention.")
	fmt.Println()
	time.Sleep(time.Second * 5)

	fmt.Println("Go do it!")
	fmt.Println("If you get distracted, come back here and type 'OK'")
	fmt.Println("If you finish, type 'Done':")

	fmt.Println()
	fmt.Println(to_do)
}

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	db_dir := filepath.Join(homeDir, ".work_motivated")
	if _, err := os.Stat(db_dir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(db_dir, 0755)
	}

	db := models.ConnectDB(db_dir)
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Are you motivated? [Y/N]")
	reader := bufio.NewReader(os.Stdin)
	user_input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	if strings.EqualFold(user_input, "n\n") {

		f := models.FeelingModel{DB: db}

		fmt.Println("Describe the joy of learning.")
		user_input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		f.Insert("joy_learning", user_input)
		fmt.Println("Read it back to yourself:")
		fmt.Println(user_input)

		time.Sleep(15 * time.Second)

		fmt.Println("Describe a time when you solved a problem.")
		user_input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		f.Insert("solve_problem", user_input)
		fmt.Println("Read it back to yourself:")
		fmt.Println(user_input)

		time.Sleep(15 * time.Second)

		fmt.Println("Describe a time when you were your best self.")
		user_input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		f.Insert("best_self", user_input)
		fmt.Println("Read it back to yourself:")
		fmt.Println(user_input)

		time.Sleep(15 * time.Second)

	}

	fmt.Println("Are you thirsty? [Y/N]")
	user_input, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}

	if strings.EqualFold(user_input, "y\n") {
		fmt.Println("Go get some water.")
	} else {
		fmt.Println("You may want to have a glass of water on your desk anyway.")
	}
	fmt.Println("Hit ENTER when you're back.")
	user_input, err = reader.ReadString('\n')

	todos := []models.ToDo{}

	fmt.Println("Check your email or notes, and list the things you want to (or should) do.")
	for {
		user_input, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.EqualFold(user_input, "\n") {
			break
		}
		todo := models.ToDo{Title: user_input}
		todos = append(todos, todo)
	}

	//What can you do pertinent to your goals/self-improvement if you get your work done quickly?
	//What can you do that you enjoy if you get your work done quickly?

	for {

		choose_one(*reader, todos)

		user_input, err = reader.ReadString('\n')

		if strings.EqualFold(user_input, "done\n") {
			fmt.Println("Sit and think.")
			fmt.Println("When you start thinking about work:")
			fmt.Println("Do 2 pushups/squats.")
			fmt.Println("Congratulate or reward yourself.")
			fmt.Println("Sit and work.")
		}

	}

}
