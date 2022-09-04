package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"github.com/soranoba/catfish/pkg/config"
	"log"
	"net/http"
	"os"
)

type (
	CmdOpts struct {
		Version        bool   `short:"v" long:"version" description:"Show the application version"`
		ConfigFilePath string `long:"config" description:"A file path of config file" required:"true"`
	}
)

func main() {
	var opts CmdOpts
	optsParser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := optsParser.Parse()

	if opts.Version {
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
		Addr:    ":8080",
		Handler: handler,
	}
	log.Fatal(srv.ListenAndServe())
}
