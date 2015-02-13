package main

import (
	"fmt"
	"reflect"
	"net/http"
	"github.com/spf13/viper"
	"github.com/gorilla/mux"
)

type Config struct {
	Redirect	string
}

func (c *Config) Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.go-url-shortener")
	viper.ReadInConfig()

	err := viper.Marshal(c)
	if (err != nil) {
		// Do proper error handling here
		fmt.Println(err)
	}
}

func (c *Config) Get(k string) (string, error) {
	// Test if the field we're requesting exists in the Config struct
	r := reflect.ValueOf(c)
	f := reflect.Indirect(r).FieldByName(k)
	if (!f.IsValid()) {
		return "", fmt.Errorf("No such field: %s in Config", k)
	}

	// All is OK! Let's return the value as a string
	return f.String(), nil
}

var c Config

func main() {
	// Instantiate the configuration
	c.Init()

	// Instantiate the mux router
	r := mux.NewRouter()
	r.HandleFunc("/s", ShortenHandler)
	r.HandleFunc("{slug:[a-zA-Z0-9]+}", ShortenedUrlHandler)
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
	// Get the redirect URL out of the Config
	redirect, err := c.Get("Redirect")
	if (err != nil) {
		// HANDLE IT!
		fmt.Println("ERRORORORORORO")
	}

	// Yay! Let's redirect ALL THE THINGS!!
	http.Redirect(w, r, redirect, 301)
}
