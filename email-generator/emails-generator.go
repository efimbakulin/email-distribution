package main

import (
	"flag"
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/email-distribution/dao"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	templatesDao *dao.TemplatesDao
	configPath   = flag.String("config", "config.cfg", "path to configuration file")
)

func main() {
	flag.Parse()

	log.Printf("Loading parameters from %s", *configPath)

	var err error
	config, err := properties.Load(*configPath)
	_ = config
	if err != nil {
		log.Fatal(err)
	}

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGTERM, syscall.SIGINT)
	s := <-stopSig
	log.Printf("Got %s - exiting", s)

	log.Print("stopped")
}
