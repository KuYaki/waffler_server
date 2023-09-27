package main

import (
	"github.com/KuYaki/waffler_server/config"
	"github.com/KuYaki/waffler_server/internal/infrastructure/logs"
	"github.com/KuYaki/waffler_server/run"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	//  load environment variables from the .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//  create configuration app
	conf := config.NewAppConf()

	// create logger
	logger, err := logs.NewLogger(conf)
	if err != nil {
		log.Fatal(err)
	}

	conf.Init(logger)

	// Creating an application instance
	app := run.NewApp(conf, logger)

	exitCode := app.
		// Initialize the application
		Bootstrap().
		// Launch the application
		Run()

	os.Exit(exitCode)
}
