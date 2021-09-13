package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type user struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"fName"`
	LastName  string `json:"lName"`
	Birthdate string `json:"bDate"`
}

type Page struct {
	Body []byte
}

type row struct {
	UsersData []user
}

var users = []user{
	{UserName: "Juan", Password: "1234", FirstName: "Juan", LastName: "Torres", Birthdate: "22"},
	{UserName: "Pablo", Password: "1234", FirstName: "Pablo", LastName: "Ramos", Birthdate: "22"},
}

func main() {
	http.HandleFunc("/index.html", viewHandler1)
	http.HandleFunc("/sign-up.html", viewHandler2)
	http.HandleFunc("/sign-in.html", viewHandler3)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/check/", checkHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func loadPage(title string) (*Page, error) {
	filename := title + ".html"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Body: body}, nil
}

func viewHandler1(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func viewHandler2(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("sign-up.html")
	t.Execute(w, p)
}

func viewHandler3(w http.ResponseWriter, r *http.Request) {
	var data *row
	data = &row{
		UsersData: users,
	}
	t, _ := template.ParseFiles("sign-in.html")

	t.Execute(w, data)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("pwd")
	firstName := r.FormValue("fName")
	lastName := r.FormValue("lName")
	bDate := r.FormValue("bDate")

	newUser := user{
		UserName:  username,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
		Birthdate: bDate}

	users = append(users, newUser)
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func checkHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username1")
	password := r.FormValue("pwd2")

	for _, a := range users {
		if a.UserName == username {
			if a.Password == password {
				http.Redirect(w, r, "/sign-in.html", http.StatusFound)
				return
			}
		}
	}
	http.Redirect(w, r, "/index.html", http.StatusFound)
}
