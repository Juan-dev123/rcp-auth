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

type table struct {
	UserName1 string
	UsersData []user
}

type message struct {
	Message string
}

var users = []user{
	{UserName: "Juan", Password: "1234", FirstName: "Juan", LastName: "Torres", Birthdate: "22"},
	{UserName: "Pablo", Password: "1234", FirstName: "Pablo", LastName: "Ramos", Birthdate: "22"},
}

var currentUser string
var messageSignIn string
var messageSignUp string

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
	messageSignUp = ""
	p := loadInfo(messageSignIn)
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, p)
}

func viewHandler2(w http.ResponseWriter, r *http.Request) {
	messageSignIn = ""
	p := loadInfo(messageSignUp)
	t, _ := template.ParseFiles("sign-up.html")
	t.Execute(w, p)
}

func viewHandler3(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var data *table
	data = &table{
		UserName1: currentUser,
		UsersData: users,
	}
	t, _ := template.ParseFiles("sign-in.html")

	t.Execute(w, data)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("pwd")
	password1 := r.FormValue("pwd1")
	firstName := r.FormValue("fName")
	lastName := r.FormValue("lName")
	bDate := r.FormValue("bDate")

	if username == "" || password == "" || firstName == "" || lastName == "" || bDate == "" {
		messageSignUp = "Please fill all the camps"
		http.Redirect(w, r, "/sign-up.html", http.StatusFound)
		return
	} else if password != password1 {
		messageSignUp = "The passwords do not coincide"
		http.Redirect(w, r, "/sign-up.html", http.StatusFound)
		return
	} else {
		newUser := user{
			UserName:  username,
			Password:  password,
			FirstName: firstName,
			LastName:  lastName,
			Birthdate: bDate}

		users = append(users, newUser)
		messageSignIn = "User registered successfully"
		http.Redirect(w, r, "/index.html", http.StatusFound)
	}
}

func checkHandler(w http.ResponseWriter, r *http.Request) {

	messageSignIn = ""
	r.ParseForm()
	username := r.FormValue("username1")
	password := r.FormValue("pwd2")

	for _, a := range users {
		if a.UserName == username {
			if a.Password == password {
				currentUser = username
				http.Redirect(w, r, "/sign-in.html", http.StatusFound)
				return
			}
		}
	}
	messageSignIn = "Username or password are incorrects"
	http.Redirect(w, r, "/index.html", http.StatusFound)
}

func loadInfo(info string) *message {
	return &message{Message: info}
}
