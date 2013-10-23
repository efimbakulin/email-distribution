package main

import (
	"bytes"
	"flag"
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/connection-string-builder"
	"github.com/efimbakulin/email-distribution/dao"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
)

var (
	consumer    *Consumer
	emailsDao   *dao.EmailsDao
	configPath  = flag.String("config", "config.cfg", "path to configuration file")
	emailRegexp = regexp.MustCompile(RE_BASIC_EMAIL)
)

const (
	RE_BASIC_EMAIL = `(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`
)

type FilterFunc func(string) bool

func IsValidEmail(email string) bool {
	return emailRegexp.MatchString(email)
}

func FilterData(input []string) []string {
	hash := make(map[string]bool)

	for _, value := range input {
		if 0 == len(value) {
			continue
		}
		if !IsValidEmail(value) {
			continue
		}

		if !hash[value] {
			hash[value] = true
		}
	}
	result := make([]string, 0, len(hash))
	for key, _ := range hash {
		result = append(result, key)
	}
	return result
}

func FilterUniqueEmails(input []byte, separator string) []string {
	bytes.Split(input, []byte(separator))
	return []string{}
}

func Handler(data []byte) error {
	emails := FilterData(strings.Split(string(data), "\n"))
	count, err := emailsDao.MarkInvalid(emails)
	if err != nil {
		log.Printf("Failed to process batch: %s", err)
		return err
	}
	log.Printf("got delivery: size: %d, processed: %d", len(emails), count)

	return nil
}

func main() {
	flag.Parse()

	var err error
	config, err := properties.Load(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	connBuilder, err := connstring.CreateBuilder(connstring.ConnectionStringPg)
	connBuilder.Address(config.String("database.addr", ""))
	connBuilder.Port(uint16(config.Int("database.port", 5432)))
	connBuilder.Username(config.String("database.username", ""))
	connBuilder.Password(config.String("database.password", ""))
	connBuilder.Dbname(config.String("database.dbname", ""))
	log.Print(connBuilder.Build())
	emailsDao = dao.NewEmailsDao(connBuilder.Build())

	if err != nil {
		log.Fatal(err)
	}

	consumer = NewConsumer(config)
	err = consumer.Connect()
	if err != nil {
		log.Fatal(err)
	}

	consumer.Serve(Handler, SkipMessageOnError)

	go func() {
		reloadSig := make(chan os.Signal, 1)
		signal.Notify(reloadSig, syscall.SIGHUP)
		s := <-reloadSig
		log.Printf("Reloading", s)
	}()

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGTERM, syscall.SIGINT)
	s := <-stopSig
	log.Printf("Got %s - exiting", s)
	consumer.Stop()
	log.Print("stopped")
}
