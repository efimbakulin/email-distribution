package main

import (
	"github.com/dmotylev/goproperties"
	"github.com/efimbakulin/connection-string-builder"
	"github.com/efimbakulin/email-distribution/email-generator/consumer/cache"
	"testing"
)

const (
	Addr     = "unstable-ebakulin.2reallife.net"
	Port     = 56432
	Username = "hv"
	Password = "oops"
	Dbname   = "sa_data"

	smtpConn = "mail.ocplay.net:25"
)

var (
	stringBuilder, _ = connstring.CreateBuilder(connstring.ConnectionStringPg)
)

func init() {
	stringBuilder.Address(Addr)
	stringBuilder.Port(Port)
	stringBuilder.Username(Username)
	stringBuilder.Password(Password)
	stringBuilder.Dbname(Dbname)

	config, _ = properties.Load("config.sample")
}

func BenchmarkGeneration(test *testing.B) {
	templates := cache.NewTemplateCache(stringBuilder.Build())
	template, err := templates.Get(2)
	if err != nil {
		test.Fatal(err)
	}

	letters := cache.NewLetterCache(stringBuilder.Build())
	letter, err := letters.Get(1)
	message := generateLetter(template, letter.Body, "test@mail.ru", 1111111)
	_ = message

	err = sendMail(letter.Subj, message, "informer-noreply@nebogame.com", "Гильдия магов", "e.v.bakulin@gmail.com")
	if err != nil {
		test.Fatal(err)
	}
}

func TestConnect(test *testing.T) {
	dialExchange(smtpConn)
}
