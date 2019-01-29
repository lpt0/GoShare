// Package config provides configuration variables
package config

import (
	"log"

	"github.com/spf13/viper"
)

// FilePath is where the uploaded files should be stored
var FilePath string

// DBPath is the path to the database file
var DBPath string

// Port is for the HTTP server
var Port string

// Authorization is a mapping of the auth-token -> uploader name
var Authorization map[string]string

// Redirects is where the server will redirect the user to if their file isn't found
var Redirects []string

// Initialize will set the above config values, using Viper
func Initialize() {
	viper.SetConfigName("config.json")
	viper.AddConfigPath(".")
	e := viper.ReadInConfig()
	if e != nil {
		log.Panicln(e)
	}
	FilePath = viper.GetString("FilePath")
	DBPath = viper.GetString("DBPath")
	Port = viper.GetString("Port")
	Authorization = viper.GetStringMapString("Authorization")
	Redirects = viper.GetStringSlice("Redirects")
}
