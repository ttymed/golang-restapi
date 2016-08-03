package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var templates = template.Must(template.ParseGlob("static/tmpl/*"))

func ShowAllTodos(w http.ResponseWriter, r *http.Request) {
	session := NewSession("mongodb://localhost")
	defer session.Close()

	results := getAllRecords("test", "todos")
	templates.ExecuteTemplate(w, "index page", nil)
	for i := 0; i < len(results); i++ {
		templates.ExecuteTemplate(w, "todo", results[i])
	}
}

func ShowSingleTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	record := getRecord("test", "todos", todoId)
	showTodoInHtml(w, record)
}

func showTodoInHtml(w http.ResponseWriter, todo Todo) {
	t := template.New("todo")
	t, _ = t.ParseFiles("static/singleTodoPage.html")
	err := t.ExecuteTemplate(w, "singleTodoPage.html", todo)
	if err != nil {
		log.Fatal(err)
	}
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
	// connect with mongoDb
	session := NewSession("mongodb://localhost")
	defer session.Close()
	c := session.DB("test").C("todos")

	r.ParseForm()
	todoId := strings.Join(r.Form["todoId"], "")
	stringTitle := strings.Join(r.Form["todoTitle"], "")
	completeDate := strings.Join(r.Form["completedBy"], "")
	record := Todo{todoId, stringTitle, false, completeDate}

	err := c.Insert(&record)
	if err != nil {
		log.Fatal(err)
	}

	ShowAllTodos(w, r)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request, todo Todo) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Println(todo)

	session := NewSession("mongodb://localhost")
	defer session.Close()

	c := session.DB("test").C("todos")
	c.Remove(bson.M{"id": todoId})

	ShowAllTodos(w, r)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

}
