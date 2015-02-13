package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	// Instantiate the configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.go-url-shortener")
	viper.ReadInConfig()

	// Instantiate the database
	var err error
	dsn := viper.GetString("mysql_user") + ":" + viper.GetString("mysql_password") + "@tcp(" + viper.GetString("mysql_host") + ":3306)/" + viper.GetString("mysql_database") + "?collation=utf8mb4_unicode_ci&parseTime=true"
	db, err = sql.Open("mysql", dsn)
	if (err != nil) {
		// HANDLE IT!
		fmt.Println(err)
	}

	defer db.Close()

	// Instantiate the mux router
	r := mux.NewRouter()
	r.HandleFunc("/s", ShortenHandler).Queries("url", "")
	r.HandleFunc("/{slug:[a-zA-Z0-9]+}", ShortenedUrlHandler)
	r.HandleFunc("/", CatchAllHandler)

	// Assign mux as the HTTP handler
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}


func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("s!"))
}

func ShortenedUrlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug, ok := vars["slug"];
	if (!ok) {
		// HANDLE IT! (a 404?)
		fmt.Println("OMGWTFBBQ")
	}
	fmt.Println(slug)

	// NOW FOR SOME DATABASE MAGIC
	var url string
	err := db.QueryRow("SELECT `url` FROM redirect WHERE slug = ?", slug).Scan(&url)
	if (err != nil) {
		// HANDLE IT! (a 404?)
		fmt.Println(err)
	}

	http.Redirect(w, r, url, 301)

	// If a route could not be matched, redirect to homepage
	// This should maybe be a 404?
	// CatchAllHandler(w, r);
}

func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	// Get the redirect URL out of the config
	if (!viper.IsSet("default_url")) {
		// HANDLE IT!
		fmt.Println("ERRORORORO")
	}

	// Yay! Let's redirect ALL THE THINGS!!
	http.Redirect(w, r, viper.GetString("default_url"), 301)
}
