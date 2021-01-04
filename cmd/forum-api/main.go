package main

import (
	"fmt"
	"forum-api/internal/app"
	"forum-api/internal/infrastructure"
	"forum-api/internal/infrastructure/persistence"
	router "forum-api/internal/interfaces/http"
	"log"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"
)

func main() {
	var srvPort string
	pflag.StringVarP(&srvPort, "port", "p", "8000", "bind port")
	var srvHost string
	pflag.StringVarP(&srvHost, "ip", "i", "", "listen addr")
	var isHelp bool
	pflag.BoolVarP(&isHelp, "help", "h", false, "usage info")
	var logLevel string
	pflag.StringVarP(&logLevel, "log-level", "l", "info", "set logging level")

	pflag.Parse()

	if isHelp {
		pflag.Usage()
		os.Exit(0)
	}

	if err := infrastructure.ConfigureLogger(os.Stdout, logLevel); err != nil {
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
		}).Fatalln("DB_PORT is not set")
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

	bdPortInt, err := strconv.Atoi(bdPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("DB_PORT is invalid")
	}

	appRepo, err := persistence.New(&persistence.DBConfig{
		Host:                 bdHost,
		Port:                 uint16(bdPortInt),
		Database:             bdName,
		User:                 bdUser,
		Password:             bdPassword,
		PreferSimpleProtocol: false,
		AcquireTimeout:       0,
		MaxConnections:       100,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot initialise DB", err)
	}

	app := app.New(appRepo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot initialise application", err)
	}

	if err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%s", srvHost, srvPort), router.New(app)); err != nil {
		logrus.WithFields(logrus.Fields{
			"pack": "main",
			"func": "main",
		}).Fatalln("Cannot start server", err)
	}
}
