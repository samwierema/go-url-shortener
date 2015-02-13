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
	// 1. Check if a slug exists
	vars := mux.Vars(r)
	slug, ok := vars["slug"];
	if (!ok) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// 2. Check if the slug exists in the database
	var url string
	err := db.QueryRow("SELECT `url` FROM `redirect` WHERE `slug` = ?", slug).Scan(&url)
	if (err != nil) {
		http.NotFound(w, r)
		return
	}

	// 3. If the slug (and thus the URL) exist, update the hit counter
	stmt, err := db.Prepare("UPDATE `redirect` SET `hits` = `hits` + 1 WHERE `slug` = ?")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = stmt.Exec(slug)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Finally, redirect the user to the URL
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get the redirect URL out of the config
	if (!viper.IsSet("default_url")) {
		// The reason for using StatusNotFound here instead of StatusInternalServerError
		// is because this is a catch-all function. You could come here via various
		// ways, so showing a StatusNotFound is friendlier than saying there's an
		// error (i.e. the configuration is missing)
		http.NotFound(w, r)
		return
	}

	// 2. If it exists, redirect the user to it
	http.Redirect(w, r, viper.GetString("default_url"), http.StatusMovedPermanently)
}
