package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish/pkg/config"
	"github.com/soranoba/henge/v2"
	"log"
	"net/http"
	"os"
)

type (
	CmdOpts struct {
		Version        bool   `short:"v" long:"version" description:"Show the application version"`
		Port           int    `short:"p" long:"port" default:"8080" description:"Bind port"`
		ConfigFilePath string `long:"config" description:"A file path of config file" required:"true"`
	}
)

func main() {
	var opts CmdOpts
	optsParser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := optsParser.Parse()

	if opts.Version {
		fmt.Println(config.AppVersion)
		os.Exit(0)
	} else if err != nil {
		log.Fatal(err)
	}

	conf, err := config.LoadYamlFile(opts.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	handler, err := NewHTTPHandler(conf)
	if err != nil {
		log.Fatal(err)
	}

	logrus.SetLevel(logrus.DebugLevel)

	srv := &http.Server{
		Addr:    ":" + henge.ToString(opts.Port),
		Handler: handler,
	}
	log.Fatal(srv.ListenAndServe())
}
