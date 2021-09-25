package main

import (
	"github.com/droplez/droplez-studio/migrations"
	"github.com/droplez/droplez-studio/pkg/server"
	"github.com/spf13/viper"
)

func init() {
	// app variables
	viper.SetDefault("environment", "dev")
	// server variables
	viper.SetDefault("droplez_studio_host", "0.0.0.0")
	viper.SetDefault("droplez_studio_port", "9090")
	// database variables
	viper.SetDefault("database_username", "droplez_studio")
	viper.SetDefault("database_password", "qwertyu9")
	viper.SetDefault("database_name", "droplez_studio")
	viper.SetDefault("database_host", "localhost")
	viper.SetDefault("database_port", "5432")
	// read environment variables that match
	viper.AutomaticEnv()
}

func main() {
	if err := migrations.Migrate(); err != nil {
		panic(err)
	}
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
