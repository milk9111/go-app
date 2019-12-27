package main

import (
	"net/http"
	"fmt"
	"html/template"
	"log"
	"os"
)

type ListOption struct {
	Name 	string
	Value 	string
	Selected bool
}

type TagCount struct {
	Label string
	Count int
}

type PageVariables struct {
	Date            string
	Time            string
	PageTitle       string
	PageListOptions []ListOption
	SubjectUrl      string
	Tag				TagCount
}

var listOptions = map[string]ListOption {
	"a": ListOption {
		Name: "Link",
		Value: "a",
		Selected: true,
	},
	"p": ListOption {
		Name: "Paragraph",
		Value: "p",
		Selected: false,
	},
	"img": ListOption {
		Name: "Image",
		Value: "img",
		Selected: false,
	},
}

func getListOptions() []ListOption {
	options := make([]ListOption, 0, len(listOptions))
	for _, value := range listOptions {
		options = append(options, value)
	}

	return options
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
	http.HandleFunc("/", DisplayPage)
	http.HandleFunc("/selected", UserSelected)
	http.ListenAndServe(port, nil)
}

func DisplayPage(w http.ResponseWriter, r *http.Request) {
	Title := "Tag Scrape"

	MyPageVariables := PageVariables{
		PageTitle: Title,
		PageListOptions: getListOptions(),
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

	url := r.Form.Get("scraping_url")
	selected := listOptions[r.Form.Get("tags")]

	updateSelected(selected)

	tag := TagCount {
		Label: selected.Name,
		Count: Scrape(url, selected.Value),
	}

	Title := "Tag Scrape"
	MyPageVariables := PageVariables{
		PageTitle: Title,
		SubjectUrl: url,
		Tag: tag,
		PageListOptions: getListOptions(),
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

func updateSelected(selectedOption ListOption) {
	for k, v := range listOptions {
		if k == selectedOption.Value {
			v.Selected = true
		} else {
			v.Selected = false
		}
	}
}
