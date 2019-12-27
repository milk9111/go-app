package main

import (
	"net/http"
	"fmt"
	"time"
	"html/template"
	"log"
	"os"
)

type RadioButton struct {
	Name 	string
	Value 	string
	IsDisabled 	bool
	IsChecked 	bool
	Text 	string
}

type PageVariables struct {
	Date 	string
	Time 	string
	PageTitle 	string
	PageRadioButtons []RadioButton
	Answer 	string
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":8080"
}

func main() {
	port := getPort()
	fmt.Println("Starting server on port", port)
	http.HandleFunc("/", DisplayRadioButtons)
	http.HandleFunc("/selected", UserSelected)
	http.ListenAndServe(port, nil)
}

func DisplayRadioButtons(w http.ResponseWriter, r *http.Request) {
	Title := "Which do you prefer?"
	MyRadioButtons := []RadioButton{
		RadioButton{"animalselect", "cats", false, false, "Cats"},
		RadioButton{"animalselect", "dogs", false, false, "Dogs"},
	}

	MyPageVariables := PageVariables{
		PageTitle: Title,
		PageRadioButtons: MyRadioButtons,
	}

	t, err := template.ParseFiles("pages/select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func UserSelected(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	youranimal := r.Form.Get("animalselect")

	Title := "Your preferred animal"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		Answer: youranimal,
	}

	t, err := template.ParseFiles("pages/select.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, MyPageVariables)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	HomePageVars := PageVariables{
		Date: now.Format("2006-02-01"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("pages/homepage.html")
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, HomePageVars)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
