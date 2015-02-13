package main

import (
	"fmt"
	"net/http"
	"github.com/spf13/viper"
	"github.com/gorilla/mux"
)

func main() {
	// Instantiate the configuration
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.go-url-shortener")
	viper.ReadInConfig()

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
	w.Write([]byte("URL!"))
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
