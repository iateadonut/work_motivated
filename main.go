package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"work/models"

	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	r         *bufio.Reader
	todos     []models.ToDo
	task_list []string
	chosen    models.ToDo
	db        *sql.DB
	sleep     func(time.Duration)
	c_timer   bool
	db_dir    string
	pprompt   promptui.Prompt
	pselect   promptui.Select
}

func sleep(d time.Duration) {
	time.Sleep(d)
}

func sleep_timer(d time.Duration) {
	time.Sleep(d)
}

func wait_for_user_to_hit_enter(app *app) {
	fmt.Println("Hit ENTER when you're back.")
	_, _ = app.r.ReadString('\n')
	return
}

func main() {

	var home_data_dir string
	var pauses bool

	flag.StringVar(&home_data_dir, "home-data-dir", ".work_motivated", "The directory within your home that will store the data.")
	flag.BoolVar(&pauses, "pauses", true, "pause for contemplation")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	db_dir := filepath.Join(homeDir, home_data_dir)
	if _, err := os.Stat(db_dir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(db_dir, 0755)
	}

	app := &app{
		pprompt: promptui.Prompt{},
		pselect: promptui.Select{},
		db_dir:  db_dir,
		todos:   []models.ToDo{},
		sleep:   sleep,
		r:       bufio.NewReader(os.Stdin),
		c_timer: false,
	}

	if false == pauses {
		app.sleep = func(d time.Duration) {}
	}

	app.db = models.ConnectDB(app.db_dir)
	if err := app.db.Ping(); err != nil {
		log.Fatal(err)
	}

	run(app)

}

func choose_one(app *app) *app {

	items := make([]string, 0, len(app.todos))
	for _, t := range app.todos {
		items = append(items, strings.Trim(t.Title, "\n"))
	}

	items = append(items, "* add another item to this list")
	items = append(items, "* finish for the day")

	prompt := promptui.Select{
		Label: "Choose one thing that you should get started on",
		Items: items,
	}

	_, to_do, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	if to_do == "* add another item to this list" {
		//fmt.Println()
		addToDos(app)
		return choose_one(app)
	}

	if to_do == "* finish for the day" {

		//What is one thing you can remember about today that will make you work better tomorrow?
		//e.g. It felt great to finish early, so make sure to get started early!

		os.Exit(1)
	}

	fmt.Println()

	for _, td := range app.todos {
		if td.Title == to_do {
			app.chosen = td
		}
	}
	return app

}

func showOneFeeling(app *app) {
	f := models.FeelingModel{DB: app.db}

	feeling := f.GetRandom()
	//fmt.Printf("%#v", feeling)

	fmt.Println()
	fmt.Println("Read your previous answer back to yourself:")
	fmt.Println(feeling.Type)
	fmt.Println(feeling.Description)
	app.sleep(time.Second * 15)
}

func motivate(app *app) {
	f := models.FeelingModel{DB: app.db}

	thinkies := map[string]string{
		"joy_learning":  "Describe the joy of learning.",
		"solve_problem": "Describe a time when you solved a problem.",
		"best_self":     "Describe a time when you were your best self.",
	}

	for field_name, instruction := range thinkies {
		fmt.Println(instruction)
		user_input, err := app.r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if !strings.EqualFold("\n", user_input) {
			f.Insert(field_name, user_input)
			fmt.Println("Read it back to yourself:")
			fmt.Println(user_input)

			app.sleep(time.Second * 15)
		}

	}

}

func timesUp(d time.Duration) {
	fmt.Println(d.String() + " is up!")
	fmt.Println("Make sure to get out of your seat and take 5 when you see this!")
	fmt.Println("Do NOT entertain yourself on your computer during your break.")
	return
	//choice: i'm finished with my break
	//i'm going to keep working without a break
}

func timer(app *app, d time.Duration) {
	if app.c_timer == false {
		app.c_timer = true
		app.sleep(d)
		timesUp(d)
		app.c_timer = false
	}
}

func anxious(app *app) error {

	//are you working on another, more important task now?

	//jot down some of the things you are distracted by.
	//things you want to look up; games you want to play; chores you think you have to do

	//save these for after work!  you may find that you don't actually want to do them, but that they were a technique for procrastination.

	//think about things you *actually* like to do; you get to do them if you finish your work

	//Repeat to yourself: I cannot fail; if I learned something I have succeeded.

	return nil
}

func bored(app *app) error {

	fmt.Println("If a particular website distracts you, 'forget it' in your history, or block it.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Tell me what you find interesting about " + app.chosen.Title)
	user_input, _ := app.r.ReadString('\n')
	fmt.Println(user_input)

	fmt.Println()
	fmt.Println("Was there anything that made you feel good when doing the last task?  Did you figure something out, or fix a problem, prevent a problem, or even just identify a problem?")
	user_input, _ = app.r.ReadString('\n')

	fmt.Println()
	fmt.Println("Let's try to see the big picture.")
	app.sleep(time.Second * 1)
	fmt.Println("Name one skill that you will acquire/improve as you do this task.")
	//book: So Good They Can't Ignore You
	user_input, _ = app.r.ReadString('\n')

	fmt.Println()
	fmt.Println("What is one result that will come of doing this task?")
	user_input, _ = app.r.ReadString('\n')

	fmt.Println()
	fmt.Println("What can you do pertinent to your goals/self-improvement if you get your work done quickly?")
	fmt.Println("Or what exciting work can you move onto when you get this work done?")
	user_input, _ = app.r.ReadString('\n')

	fmt.Println()
	fmt.Println("What can you do later that you enjoy if you focus now on this task?")
	user_input, _ = app.r.ReadString('\n')

	fmt.Println()
	fmt.Println("Repeat to yourself:")
	quotes := []string{
		"After working undistracted for some time, I may acheive a flow state, with all of its benefits",
		"My sense of boredom will pass after just a few minutes of working on " + app.chosen.Title,
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(quotes))
	quote := quotes[randomIndex]
	fmt.Println(quote)

	app.sleep(time.Second * 6)

	return nil
}

func icky(app *app) error {

	fmt.Println("Clean your desk.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Brush your teeth or otherwise tend to your oral health.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Trim your facial hair if you need and are able.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Put on chap stick if needed.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	app.pprompt.Label = "Do you need to take a shower? [Y/N]"
	user_input, err := app.pprompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	if strings.EqualFold(user_input, "y") {
		fmt.Println("Take a quick cold shower if you can.")
		fmt.Println("Or remember to schedule a shower asap.")
		fmt.Println("Hit ENTER when complete")
		_, _ = app.r.ReadString('\n')
	}

	return nil
}

func overwhelmed(app *app) error {

	fmt.Println("Clean your desk.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Close *all* windows & tabs not related to the task at hand.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	//do you think daily, taking time for yourself with no tv, handphone, etc, just to sit down, or lie down, or pace and think for at least 10 minutes?

	fmt.Println("If a particular website distracts you, 'forget it' in your history, or block it.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Turn off any background noise.")
	fmt.Println("If you need earplugs, put them in.")
	fmt.Println("If you really think you need background noise, put on some 'white noise' or NON-lyrical music. Do not put on any background noise with any words at all.")
	fmt.Println("Hit ENTER when complete")
	_, _ = app.r.ReadString('\n')

	//box breathing

	// fmt.Println()
	// fmt.Println("Sit and think.")
	// fmt.Println("When you start thinking about work:")
	// fmt.Println("Do 2 pushups/squats.")
	// fmt.Println("Sit and work.")
	// fmt.Println()

	//intrusive thoughts

	return nil
}

func tired(app *app) error {

	app.pprompt.Label = "Are you able to take a siesta/power nap? [Y/N]"
	user_input, err := app.pprompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return err
	}

	if strings.EqualFold(user_input, "y") {
		fmt.Println("Remind yourself:")
		fmt.Println("If I complete my tasks quickly, I will be able to lie down for a nap.")
		fmt.Println()
		app.sleep(time.Second * 5)
	}

	fmt.Println("Drink some water.")
	fmt.Println("Hit ENTER.")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Get up; do some squats or pushups, not too many.")
	fmt.Println("Hit ENTER when you're back.")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Repeat to yourself:")
	fmt.Println("When I am tired and work with focus, I find that work invigorates me after a little while.")
	fmt.Println()
	app.sleep(time.Second * 5)

	return nil
}

func distracted(app *app) error {

	//are you still working on smallest_task?

	fmt.Println("What is the smallest action you can take right now to get going?")

	smallest_task, err := app.r.ReadString('\n')
	if err != nil {
		return err
	}
	//fmt.Println(smallest_task)

	app.chosen.Smalltask = strings.Trim(smallest_task, "\r\n")

	if false == app.c_timer {
		fmt.Println("Set a timer for 25 minutes.  Then hit enter.")
		_, _ = app.r.ReadString('\n')

		fmt.Println("What can you reasonably get done in 25 minutes?")
		fmt.Println("Hit ENTER on an empty line to finish.")
		app.task_list = nil
		for {
			user_input, err := app.r.ReadString('\n')
			if err != nil {
				fmt.Println(err)
			}
			app.task_list = append(app.task_list, strings.Trim(user_input, "\r\n"))
			if strings.EqualFold(user_input, "\n") {
				break
			}
			if strings.EqualFold(user_input, "\r\n") {
				break
			}
		}

		//fmt.Println("Try to work for the duration of the timer without distraction.")
		fmt.Println("Do your best to finish your list before the timer ends.")
		go timer(app, time.Duration(time.Minute*25))
	} else {
		prompt := promptui.Select{
			Label: "How do you feel?",
			Items: []string{"anxious", "bored", "icky", "overwhelmed", "tired"},
		}

		_, reason, err := prompt.Run()
		if err != nil {
			return err
		}

		switch reason {
		case "anxious":
			anxious(app)

		case "bored":
			bored(app)

		case "icky":
			icky(app)

		case "overwhelmed":
			overwhelmed(app)

		case "tired":
			tired(app)
		}

	}

	fmt.Println()

	// fmt.Println(app.chosen.Title)
	// fmt.Println(app.chosen.Smalltask)
	// fmt.Println()

	//Are you working by hours or tasks?
	//if hours
	//that's not great.  if you have to work by hours, you tend to stretch work out to fill those hours.  we'll work on that in a second, for now:
	//How many hours do you need to work before a longer break?
	//if you start right now and work straight through, you can finish at x p.m.

	//How many tasks can you reasonably get done in that amount of time?

	//What task do you need to get done?

	// fmt.Println("Type 'done' if you complete '" + app.chosen.Title + "'; anything else if distracted:")

	// []string{"done", "distracted", "pause for a break"}

	return nil

}

func addToDos(app *app) {
	fmt.Println("Add them line by line. Hit ENTER on an empty line to finish.")
	for {
		user_input, err := app.r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		if strings.EqualFold(user_input, "\n") {
			break
		}
		if strings.EqualFold(user_input, "\r\n") {
			break
		}
		todo := models.ToDo{Title: strings.Trim(user_input, "\r\n")}
		app.todos = append(app.todos, todo)
	}
	return
}

func goDoIt(app *app) (string, error) {
	fmt.Println("Repeat to yourself:")
	//pick a random phrase
	intentions := []string{
		"I work with intention.",
		"I am proud of my work.",
		"I do all things purposefully.",
	}
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(intentions))
	intention := intentions[randomIndex]

	fmt.Println(intention)
	fmt.Println()
	app.sleep(time.Second * 5)

	// fmt.Println("Go do it!")
	// app.sleep(time.Second * 1)

	// fmt.Println("When you finish, select 'Finished'.")
	// fmt.Println("If you get distracted, select 'Distracted'.")

	fmt.Println()
	//fmt.Println(to_do)

	app.pselect.Label = "How is it going with '" + app.chosen.Title + "'?"
	app.pselect.Items = []string{"Finished", "Distracted", "Change Tasks"}
	_, user_input, err := app.pselect.Run()
	return user_input, err
}

func run(app *app) {

	app.pprompt.Label = "Are you motivated? [Y/N]"
	user_input, err := app.pprompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	// fmt.Println()

	if strings.EqualFold(user_input, "n") {
		motivate(app)
	} else {
		showOneFeeling(app)
	}

	app.pprompt.Label = "Are you thirsty? [Y/N]"
	user_input, err = app.pprompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	// fmt.Println()

	if strings.EqualFold(user_input, "y") {
		fmt.Println("Go get some water.")
	} else {
		fmt.Println("You may want to have a glass of water on your desk anyway.")
	}
	fmt.Println("Hit ENTER when you're back.")
	_, _ = app.r.ReadString('\n')

	fmt.Println("Check your email or notes, and list the things you want to (or should) do.")

	addToDos(app)

	for {

		//fmt.Print("\a")
		choose_one(app)
		//fmt.Printf("%#v", app.chosen)

	InputLoop:
		for {
			time_start := time.Now()

			user_input, err = goDoIt(app)

			time_taken := time.Since(time_start).Round(time.Minute).Minutes()
			fmt.Println("You've been working " + strconv.FormatFloat(time_taken, 'f', 0, 64) + " minutes")
			app.sleep(time.Second * 3)

			if strings.EqualFold(user_input, "finished") {

				// fmt.Printf("%#v", app.chosen)

				for idx, t := range app.todos {
					if t.Title == app.chosen.Title {
						app.todos = append(app.todos[0:idx], app.todos[idx+1:]...)
					}
				}
				fmt.Println("Congratulate yourself.")
				fmt.Println()
				app.sleep(time.Second * 7)
				break InputLoop

			} else if strings.EqualFold(user_input, "change tasks") {
				choose_one(app)
			} else if strings.EqualFold(user_input, "distracted") {
				distracted(app)
			}
		}

	}
}

//This function only demonstrates how to mock multiple inputs on a single function.  It is used in TestMock().
func mock(p promptui.Prompt) string {
	p.Label = "[Y/N]"
	user_input, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	user_input2, err := p.Run()

	return user_input + user_input2
}
