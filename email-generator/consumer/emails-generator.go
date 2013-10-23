package main

import (
	"bytes"
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/connection-string-builder"
	"github.com/efimbakulin/email-distribution/email-generator/consumer/cache"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"
)

var (
	configPath = flag.String("config", "config.sample", "path to configuration file")
	letters    *cache.Letters
	templates  *cache.Templates
	config     properties.Properties
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

func main() {
	flag.Parse()

	var err error
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

	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGTERM, syscall.SIGINT)
	s := <-stopSig
	log.Printf("Got %s - exiting", s)
}
