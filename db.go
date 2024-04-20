package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Article struct {
	ID      int
	Title   string
	Content string
}

func Connect() {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	log.Println("datbaase opened")
	db.Exec(
		"CREATE TABLE IF NOT EXISTS posts (id integer not null primary key autoincrement, title text, content text);",
	)
	DB = db
	log.Println("Connected to DB")
}

func GetArticles() []Article {
	query, err := DB.Prepare("SELECT * FROM posts")
	if err != nil {
		panic(err)
	}
	rows, err := query.Query()
	if err != nil {
		panic(err)
	}
	var articles []Article
	for rows.Next() {
		var article Article
		rows.Scan(&article.ID, &article.Title, &article.Content)
		articles = append(articles, article)
	}
	return articles
}

func GetArticleByID(id int) Article {
	query, err := DB.Prepare("SELECT * FROM posts WHERE id=?")
	if err != nil {
		panic(err)
	}
	rows, err := query.Query(id)
	var article Article
	for rows.Next() {
		rows.Scan(&article.ID, &article.Title, &article.Content)
	}
	return article
}

func DeleteArticle(id int) {
	query, err := DB.Prepare("DELETE FROM posts WHERE id=?")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(id)
	if err != nil {
		panic(err)
	}
}

func CreateArticle(article Article) {
	query, err := DB.Prepare("INSERT INTO posts (title,content) VALUES (?,?)")
	if err != nil {
		panic(err)
	}
	_, err = query.Exec(article.Title, article.Content)
	if err != nil {
		panic(err)
	}
}
