package main

import (
	"encoding/json"
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

type table struct {
	UserName1 string
	UsersData []user
}

type message struct {
	Message string
}

var users []user

var currentUser string
var messageSignIn string
var messageSignUp string

const DataPath string = "data/users.txt"

func main() {
	readFile()
	http.HandleFunc("/index.html", viewHandler1)
	http.HandleFunc("/sign-up.html", viewHandler2)
	http.HandleFunc("/sign-in.html", viewHandler3)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/check/", checkHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

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
		writeFile()
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

func writeFile() {
	content, err1 := json.Marshal(&users)
	if err1 != nil {
		log.Fatal(err1)
	} else {
		err := ioutil.WriteFile(DataPath, content, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func readFile() {
	content, err := ioutil.ReadFile(DataPath)
	if err != nil {
		log.Fatal(err)
	} else {
		json.Unmarshal(content, &users)
	}
}
