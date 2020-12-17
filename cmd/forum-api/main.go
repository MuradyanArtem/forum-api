package main

import (
	"fmt"
	"forum-api/internal/app"
	"forum-api/internal/pkg/infrastructure/persistence"
	"forum-api/internal/pkg/infrastructure/utils"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

func main() {
	var srvPort string
	pflag.StringVarP(&srvPort, "port", "p", "8000", "bind port")
	var srvHost string
	pflag.StringVarP(&srvHost, "host", "ip", "", "host addr")

	var isHelp bool
	pflag.BoolVarP(&isHelp, "help", "h", false, "usage info")
	var logLevel string
	pflag.StringVarP(&logLevel, "log-level", "l", "info", "set logging level")

	pflag.Parse()

	if isHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if err := utils.ConfigureLogger(os.Stdout, logLevel); err != nil {
		log.Fatalln(err, "Cannot initialise logger")
	}

	bdHost, exist := os.LookupEnv("DB_HOST")
	if !exist {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_HOST is not set")
	}
	bdPort, exist := os.LookupEnv("DB_PORT")
	if !exist {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_HOST is not set")
	}
	bdName, exist := os.LookupEnv("DB_NAME")
	if !exist {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_NAME is not set")
	}
	bdUser, exist := os.LookupEnv("DB_USER")
	if !exist {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_USER is not set")
	}
	bdPassword, exist := os.LookupEnv("DB_PASSWORD")
	if !exist {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_PASSWORD is not set")
	}

	appRepo, err := persistence.New(&persistence.DBConfig{
		Host: bdHost,
		Port: bdPort,
		Database: bdName,
		User:     bdUser,
		Password: bdPassword,
		PreferSimpleProtocol: false,
		AcquireTimeout: 0,
		MaxConnections: 100,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot initialise DB", err)
	}

	app, err := app.New(appRepo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot initialise application", err)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", srvHost, srvPort),
		Handler:      router.New(app),
		WriteTimeout: time.Duration(3) * time.Second,
		ReadTimeout:  time.Duration(3) * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot start server", err)
	}
}
