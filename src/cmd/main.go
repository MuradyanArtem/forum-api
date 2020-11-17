package main

import (
	"os"
	"flag"
	"log"

	"gopkg.in/yaml.v2"

	"apiserver/src/internal/app/apiserver"
)

var (
	configPath string,
	isHelp bool
)

func init() {
	flag.StringVar(&configPath, "-f", "../configs/production.yaml", "path to config file")
	flag.BoolVar(&isHelp, "-h", false, "show usage info")
}

func main() {
	flag.Parse()
	if isHelp {
		flag.Usage()
		os.Exit(0)
	}

	var srvConf apiserver.Config
	if err := apiserver.ParseConfig(configPath, srvConf) != nil {
		log.Fatal(err)
	}

	if err := apiserver.Run(srvConf); err != nil {
		log.Fatal(err)
	}
}
