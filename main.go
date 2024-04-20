package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var count = 0

func Render(templateName string, w http.ResponseWriter, data interface{}) {
	tmpl, err := template.ParseFiles("templates/base.html", "templates/"+templateName)
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

func loadPosts(w http.ResponseWriter, _ *http.Request) {
	posts := GetArticles()
	var html string = "<ul>"
	for _, post := range posts {
		html += "<li>"
		html += "<a href=/posts/" + strconv.Itoa(post.ID) + "><h2>" + post.Title + "</h2>"

		html += "</li>"

	}
	html += "</ul>"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"count": count,
	}
	Render("home.html", w, data)
}

func PostDetials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	post := GetArticleByID(id)
	t, err := template.ParseFiles("templates/post_details.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, post)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	article := Article{Title: r.Form.Get("title"), Content: r.Form.Get("content")}
	CreateArticle(article)
	articles := GetArticles()
	html := "<ul>"
	for _, article := range articles {
		html += "<li>"
		html += "<h2>" + article.Title + "<h2/>"
		html += "<p>" + article.Content + "<p/>"
		html += "</li>"
	}
	html += "</ul>"
	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(html))
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	DeleteArticle(id)
	w.Write([]byte("<h1>Post Deleted</h1><a href=/>Return Home</a>"))
}

func main() {
	Connect()

	r := mux.NewRouter()
	r.HandleFunc("/", HomePage)
	r.HandleFunc("/create", createPost).Methods("POST")
	r.HandleFunc("/posts/{id}", PostDetials).Methods("GET")
	r.HandleFunc("/delete/{id}", DeletePost).Methods("DELETE")
	r.HandleFunc("/load", loadPosts).Methods("GET")
	log.Println("Server listening")
	log.Println(http.ListenAndServe(":5000", r))
}
