package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/connection-string-builder"
	"github.com/efimbakulin/email-distribution/dao"
	"github.com/efimbakulin/email-distribution/email-generator/consumer/cache"
	"io"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

var (
	consumer        *Consumer
	configPath      = flag.String("config", "config.sample", "path to configuration file")
	letters         *cache.Letters
	templates       *cache.Templates
	emails          *dao.Emails
	config          properties.Properties
	applicationName string
)

func calcMd5(val string) string {
	h := md5.New()
	io.WriteString(h, val)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func generateUnsubscribtionLink(secret string, pageUrl string, email string, emailId int64) string {
	hash := generateUnsubscribtionHash(secret, email, emailId)
	return fmt.Sprintf("%semail=%s&hash=%s", pageUrl, email, hash)
}

func generateUnsubscribtionHash(secret string, email string, emailId int64) string {
	return calcMd5(calcMd5(fmt.Sprintf("%s%s%s%s", emailId, email, secret)) + secret)
}

func generateLetter(tpl *template.Template, body string, email string, emailId int64) string {
	params := make(map[string]string)
	params["Body"] = body
	params["UnsubscribeLink"] = generateUnsubscribtionLink(
		config.String("unsubscribtion.secret_key", ""),
		config.String("unsubscribtion.page_url", ""),
		email,
		emailId,
	)
	var result bytes.Buffer
	tpl.Execute(&result, params)
	return result.String()
}

func startConsumer() error {
	consumer = NewConsumer(config)
	err := consumer.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Rabbitmq listener started")
	go consumer.Serve(AmqpHandler, SkipMessageOnError)
	return nil
}

type IncomeMessage struct {
	Body     string
	To       string
	From     string
	FromName string
	Language string
}

func AmqpHandler(data []byte) error {
	log.Print("New data got")
	var js IncomeMessage
	if err := json.Unmarshal(data, &js); err != nil {
		log.Printf("%s", err)
		return err
	}
	id, err := emails.GetId(js.To)
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	_ = id
	log.Printf("%s", js)
	return fmt.Errorf("Not ready yet")
}

func main() {
	flag.Parse()

	w, err := syslog.New(syslog.LOG_INFO, applicationName)
	if err != nil {
		log.Fatalf("connecting to syslog: %s", err)
	}

	log.SetOutput(w)
	log.SetFlags(0)

	config, err = properties.Load(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	connBuilder, err := connstring.CreateBuilder(connstring.ConnectionStringPg)
	connBuilder.Address(config.String("database.addr", ""))
	connBuilder.Port(uint16(config.Int("database.port", 5432)))
	connBuilder.Username(config.String("database.username", ""))
	connBuilder.Password(config.String("database.password", ""))
	connBuilder.Dbname(config.String("database.dbname", ""))

	templates = cache.NewTemplateCache(connBuilder.Build())
	letters = cache.NewLetterCache(connBuilder.Build())
	emails = dao.NewEmailsDao(connBuilder.Build())

	if err = startConsumer(); err != nil {
		log.Fatal(err)
	}

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGTERM, syscall.SIGINT)
	s := <-stopSig
	log.Printf("Got %s - exiting", s)
}
